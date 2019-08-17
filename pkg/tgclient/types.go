package tgclient

import (
	"encoding/json"
)

type Request map[string]interface{}

func (d Request) String() string {
	bytes, _ := json.Marshal(d)
	return string(bytes)
}

type typeHolder struct {
	Type ClassType `json:"@type"`
}

type ClassType string

const (
	ErrorEventType       ClassType = "error"
	NewMessageUpdateType ClassType = "updateNewMessage"
	MessageTextType      ClassType = "messageText"
)

type rawEvent struct {
	Type  ClassType `json:"@type,omitempty"`
	Extra string    `json:"@extra,omitempty"`
}

type Event struct {
	Type     ClassType       `json:"@type"`
	Extra    string          `json:"@extra"`
	Contents json.RawMessage `json:"contents"`
}

func (e Event) Unmarshal(obj interface{}) (err error) {
	err = json.Unmarshal(e.Contents, obj)
	if err != nil {
		err = ParseErr.Wrap(err, "unmarshal event failed. ev: %s", e.String())
	}
	return
}

func (e Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

type ErrorEvent struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NewMessageUpdate struct {
	Message Message `json:"message"`
}

type AuthState string

const (
	AuthStateClosed              AuthState = "authorizationStateClosed"
	AuthStateClosing             AuthState = "authorizationStateClosing"
	AuthStateLoggingOut          AuthState = "authorizationStateLoggingOut"
	AuthStateReady               AuthState = "authorizationStateReady"
	AuthStateWaitCode            AuthState = "authorizationStateWaitCode"
	AuthStateWaitEncryptionKey   AuthState = "authorizationStateWaitEncryptionKey"
	AuthStateWaitPassword        AuthState = "authorizationStateWaitPassword"
	AuthStateWaitPhoneNumber     AuthState = "authorizationStateWaitPhoneNumber"
	AuthStateWaitTdlibParameters AuthState = "authorizationStateWaitTdlibParameters"
)

type Chat struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

func (c Chat) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

type User struct {
	Id        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

func (u User) String() string {
	bytes, _ := json.Marshal(u)
	return string(bytes)
}

type Messages struct {
	TotalCount int       `json:"total_count"`
	Messages   []Message `json:"messages"`
}

type Message struct {
	Id           int64           `json:"id"`
	ChatId       int64           `json:"chat_id"`
	SenderUserId int32           `json:"sender_user_id"`
	IsOutgoing   bool            `json:"is_outgoing"`
	RawContent   json.RawMessage `json:"content"`
}

func (m Message) UnmarshalContent(obj interface{}) error {
	err := json.Unmarshal(m.RawContent, obj)
	if err != nil {
		return ParseErr.WrapWithNoMessage(err)
	}
	return nil
}

func (m Message) GetContentType() (ClassType, error) {
	h := typeHolder{}
	err := json.Unmarshal(m.RawContent, &h)
	if err != nil {
		return "", ParseErr.WrapWithNoMessage(err)
	}
	return h.Type, nil
}

func (m Message) String() string {
	raw, _ := json.Marshal(m)
	return string(raw)
}

type MessageText struct {
	Text FormattedText `json:"text"`
}

type FormattedText struct {
	Text string `json:"text"`
}
