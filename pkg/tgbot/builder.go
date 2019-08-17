package tgbot

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"time"
)

type socks5Proxy struct {
	host     string
	port     int
	login    string
	password string
}


type Builder struct {
	token string
	proxy *socks5Proxy
	timeout time.Duration
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Socks5Proxy(host string, port int, login, password string) *Builder {
	b.proxy = &socks5Proxy{
		host:     host,
		port:     port,
		login:    login,
		password: password,
	}
	return b
}

func (b *Builder) Token(token string) *Builder {
	b.token = token
	return b
}

func (b *Builder) TimeoutSec(secs int) *Builder {
	b.timeout = time.Duration(secs) * time.Second
	return b
}

func (b *Builder) Build() (*Bot, error) {
	client := http.Client{
		Timeout: b.timeout,
	}

	if b.proxy != nil {
		dialer, err := b.prepareProxy()
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}

	return &Bot{
		token: b.token,
		client: &client,

	}, nil
}

func (b *Builder) prepareProxy() (d proxy.Dialer, err error) {
	addr := fmt.Sprintf("%s:%d", b.proxy.host, b.proxy.port)
	auth := proxy.Auth{
		User: b.proxy.login,
		Password: b.proxy.password,
	}
	d, err = proxy.SOCKS5("tcp", addr, &auth, proxy.Direct)
	if err != nil {
		err = BuilderErr.Wrap(err, "proxy connect failed")
	}
	return
}
