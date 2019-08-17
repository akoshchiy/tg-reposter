package tgclient

/*
#cgo LDFLAGS: -ltdjson
#include <stdlib.h>
#include <td/telegram/td_json_client.h>
#include <td/telegram/td_log.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var ReceiveTimeout = 10.0
var RequestTimeout = time.Second * 300000


type Client struct {
	logger     *logrus.Entry
	config     config
	client     unsafe.Pointer
	idGen      uint64
	closed     int32
	proxyAdded int32

	reqMu    sync.Mutex
	requests map[uint64]chan Event

	eventsMu sync.Mutex
	events   map[ClassType]chan Event
}

func newClient(config config) *Client {
	client := &Client{
		logger:   logrus.WithField("logger", "tgclient"),
		config:   config,
		client:   C.td_json_client_create(),
		reqMu:    sync.Mutex{},
		requests: map[uint64]chan Event{},
		eventsMu: sync.Mutex{},
		events:   map[ClassType]chan Event{},
	}
	go client.updateLoop()
	return client
}

func (c *Client) Destroy() {
	C.td_json_client_destroy(c.client)
	c.setClosed()
}

func (c *Client) Send(r Request) (Event, error) {
	id := atomic.AddUint64(&c.idGen, 1)

	req := C.CString(c.prepareRequest(id, r))
	defer C.free(unsafe.Pointer(req))

	wait := c.newWaitChan(id)

	C.td_json_client_send(c.client, req)

	resp, err := c.waitResponse(r, wait)
	c.removeWaitChan(id)

	return resp, err
}

func (c *Client) SendAndForget(r Request) {
	req := C.CString(c.prepareRequest(0, r))
	defer C.free(unsafe.Pointer(req))
	C.td_json_client_send(c.client, req)
}

func (c *Client) waitResponse(req Request, ch chan Event) (Event, error) {
	select {
	case resp := <-ch:
		if resp.Type == ErrorEventType {
			return Event{}, c.handleError(req, resp)
		}
		return resp, nil
	case <-time.After(RequestTimeout):
		return Event{}, TimeoutErr.New("req timeout: " + req.String())
	}
}

func (c *Client) handleError(req Request, raw Event) error {
	errorEv := ErrorEvent{}
	err := raw.Unmarshal(&errorEv)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("req failed. code: %d, msg: %s, reg: %s", errorEv.Code, errorEv.Message, req.String())
	return RequestErr.New(msg)
}

func (c *Client) prepareRequest(id uint64, data Request) string {
	data["@extra"] = strconv.FormatUint(id, 10)
	req, _ := json.Marshal(data)
	return string(req)
}

func (c *Client) newWaitChan(id uint64) chan Event {
	c.reqMu.Lock()
	defer c.reqMu.Unlock()
	wait := make(chan Event)
	c.requests[id] = wait
	return wait
}

func (c *Client) removeWaitChan(id uint64) {
	c.reqMu.Lock()
	defer c.reqMu.Unlock()
	delete(c.requests, id)
}

func (c *Client) updateLoop() {
	for !c.checkClosed() {
		event, err := c.receive(ReceiveTimeout)
		if err != nil {
			c.logger.Errorf("receive error: %+v", err)
			continue
		}
		if event != nil {
			c.handleEvent(*event)
		}
	}
}

func (c *Client) checkClosed() bool {
	v := atomic.LoadInt32(&c.closed)
	return v == 1
}

func (c *Client) setClosed() {
	atomic.StoreInt32(&c.closed, 1)
}

func (c *Client) handleEvent(ev Event) {
	if ev.Extra == "" {
		c.fireEvent(ev)
		return
	}
	id, err := strconv.ParseUint(ev.Extra, 10, 64)
	if err != nil {
		c.logger.Errorf("failed to parse extra. update: %s", ev.String())
		return
	}
	c.handleResponse(id, ev)
}

func (c *Client) fireEvent(ev Event) {
	c.logger.Tracef("event: %s", ev)

	c.eventsMu.Lock()
	defer c.eventsMu.Unlock()

	if ch, ok := c.events[ev.Type]; ok {
		ch <- ev
	}
}

func (c *Client) handleResponse(id uint64, ev Event) {
	c.reqMu.Lock()
	defer c.reqMu.Unlock()

	if ch, ok := c.requests[id]; ok {
		ch <- ev
	}
}

func (c *Client) receive(timeout float64) (*Event, error) {
	resp := C.td_json_client_receive(c.client, C.double(timeout))
	if resp == nil {
		return nil, nil
	}

	raw := rawEvent{}
	contents := json.RawMessage([]byte(C.GoString(resp)))

	err := json.Unmarshal(contents, &raw)
	if err != nil {
		return nil, ParseErr.Wrap(err, "failed to parse receive data")
	}

	return &Event{
		Type:     raw.Type,
		Extra:    raw.Extra,
		Contents: contents,
	}, nil
}

func (c *Client) isProxyAdded() bool {
	v := atomic.LoadInt32(&c.proxyAdded)
	return v == 1
}

func (c *Client) setProxyAdded(val bool) {
	var added int32 = 0
	if val {
		added = 1
	}
	atomic.StoreInt32(&c.proxyAdded, added)
}
