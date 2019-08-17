package tgbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiUrl = "https://api.telegram.org/bot%s/%s"

type request map[string]interface{}

func (r request) String() string {
	raw, _ := json.Marshal(r)
	return string(raw)
}

type response struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
	ErrorCode   int             `json:"error_code"`
}

type Bot struct {
	token  string
	client *http.Client
}

func (b *Bot) GetMe() (u User, err error) {
	resp, err := b.doRequest("getMe", request{})
	if err != nil {
		return
	}
	err = json.Unmarshal(resp.Result, &u)
	if err != nil {
		err = ReqErr.WrapWithNoMessage(err)
	}
	return
}

func (b *Bot) SendMessage(chatId int64, text string) (err error) {
	req := request{
		"chat_id": chatId,
		"text": text,
	}
	_, err = b.doRequest("sendMessage", req)
	return
}

func (b *Bot) ForwardMessage(chatId, fromChatId, messageId int64) (err error) {
	req := request{
		"chat_id":      chatId,
		"from_chat_id": fromChatId,
		"message_id":   messageId,
	}
	_, err = b.doRequest("forwardMessage", req)
	return
}

func (b *Bot) doRequest(method string, req request) (resp response, err error) {
	jsonStr, _ := json.Marshal(req)
	url := b.getUrl(method)
	buf := bytes.NewBuffer(jsonStr)

	httpResp, err := b.client.Post(url, "application/json", buf)

	if httpResp != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		err = newReqError(err, method, req)
		return
	}
	err = json.NewDecoder(httpResp.Body).Decode(&resp);
	if err != nil {
		err = newReqError(err, method, req)
	}
	if !resp.Ok {
		err = newApiError(err, method, req, resp)
	}
	return
}

func (b *Bot) getUrl(method string) string {
	return fmt.Sprintf(apiUrl, b.token, method)
}
