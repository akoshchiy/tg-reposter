package tgclient

func (c *Client) GetAuthState() (AuthState, error) {
	resp, err := c.Send(Request{"@type": "getAuthorizationState"})
	if err != nil {
		return "", err
	}
	return AuthState(resp.Type), nil
}

func (c *Client) SetLogVerbosity(verbosity int) {
	c.SendAndForget(Request{
		"@type":               "setLogVerbosityLevel",
		"new_verbosity_level": verbosity,
	})
}

func (c *Client) Authorize() error {
	state, err := c.GetAuthState()
	if err != nil {
		return err
	}

	if state == AuthStateWaitTdlibParameters {
		err = c.setTdLibParameters()
		if err != nil {
			return err
		}
		return c.Authorize()
	}

	if c.config.proxy != nil && !c.isProxyAdded() {
		c.addProxy()
		c.setProxyAdded(true)
	}

	if state == AuthStateWaitEncryptionKey {
		err = c.checkDatabaseEncryptionKey(nil)
		if err != nil {
			return err
		}
		return c.Authorize()
	}

	if state == AuthStateWaitPhoneNumber {
		err = c.setAuthenticationPhoneNumber()
		if err != nil {
			return err
		}
		return c.Authorize()
	}

	if state == AuthStateWaitCode {
		err = c.checkAuthenticationCode()
		if err != nil {
			return err
		}
		return c.Authorize()
	}

	if state == AuthStateWaitPassword {
		err = c.checkAuthenticationPassword()
		if err != nil {
			return err
		}
		return c.Authorize()
	}

	if state == AuthStateReady {
		return nil
	}
	return AuthErr.New("auth failed. state: " + string(state))
}

func (c *Client) GetMe() (u User, err error) {
	r := Request{"@type": "getMe"}
	ev, err := c.Send(r)
	if err != nil {
		return
	}
	err = parseResponse(ev, r, &u)
	return
}

func (c *Client) GetChats(offsetOrder, offsetChatId, limit int64) ([]int64, error) {
	type rawChats struct {
		ChatIds []int64 `json:"chat_ids"`
	}
	r := Request{
		"@type":          "getChats",
		"offset_order":   offsetOrder,
		"offset_chat_id": offsetChatId,
		"limit":          limit,
	}
	ev, err := c.Send(r)
	if err != nil {
		return nil, err
	}
	chats := rawChats{}
	err = parseResponse(ev, r, &chats)
	return chats.ChatIds, err
}

func (c *Client) GetChat(chatId int64) (ch Chat, err error) {
	r := Request{
		"@type":   "getChat",
		"chat_id": chatId,
	}
	ev, err := c.Send(r)
	if err != nil {
		return
	}
	err = parseResponse(ev, r, &ch)
	return
}

func (c *Client) GetChatHistory(chatId int64, fromMsgId int64, offset int, limit int) (m Messages, err error) {
	r := Request{
		"@type":           "getChatHistory",
		"chat_id":         chatId,
		"from_message_id": fromMsgId,
		"offset":          offset,
		"limit":           limit,
	}
	ev, err := c.Send(r)
	if err != nil {
		return
	}
	err = parseResponse(ev, r, &m)
	return
}

func (c *Client) ListenNewMessages() <-chan Message {
	eventCh := c.addEventChannel(NewMessageUpdateType)

	ch := make(chan Message)

	go func() {
		for ev := range eventCh {
			update := NewMessageUpdate{}
			err := ev.Unmarshal(&update)
			if err != nil {
				c.logger.Errorf("%+v", err)
			} else {
				ch <- update.Message
			}
		}
	}()

	return ch
}

func (c *Client) addEventChannel(t ClassType) chan Event {
	c.eventsMu.Lock()
	defer c.eventsMu.Unlock()

	_, ok := c.events[NewMessageUpdateType]
	if ok {
		panic("already listening event: " + NewMessageUpdateType)
	}

	ch := make(chan Event)
	c.events[NewMessageUpdateType] = ch

	return ch
}

func (c *Client) checkAuthenticationPassword() error {
	data := Request{
		"@type":    "checkAuthenticationPassword",
		"password": c.config.password,
	}
	_, err := c.Send(data)
	return err
}

func (c *Client) checkAuthenticationCode() error {
	data := Request{
		"@type":      "checkAuthenticationCode",
		"code":       c.config.checkCode,
		"first_name": "",
		"last_name":  "",
	}
	_, err := c.Send(data)
	return err
}

func (c *Client) setAuthenticationPhoneNumber() error {
	data := Request{
		"@type":                   "setAuthenticationPhoneNumber",
		"phone_number":            c.config.authPhone,
		"allow_flash_call":        false,
		"is_current_phone_number": false,
	}
	_, err := c.Send(data)
	return err
}

func (c *Client) checkDatabaseEncryptionKey(key []byte) error {
	data := Request{
		"@type":          "checkDatabaseEncryptionKey",
		"encryption_key": key,
	}
	_, err := c.Send(data)
	return err
}

func (c *Client) setTdLibParameters() error {
	data := Request{
		"@type": "setTdlibParameters",
		"parameters": Request{
			"@type":                    "tdlibParameters",
			"database_directory":       c.config.databaseDirectory,
			"use_test_dc":              false,
			"files_directory":          c.config.filesDirectory,
			"use_file_database":        c.config.useFileDatabase,
			"use_chat_info_database":   false,
			"use_message_database":     false,
			"use_secret_chats":         false,
			"api_id":                   c.config.apiId,
			"api_hash":                 c.config.apiHash,
			"system_language_code":     c.config.systemLanguageCode,
			"device_model":             c.config.deviceModel,
			"system_version":           c.config.systemVersion,
			"application_version":      c.config.applicationVersion,
			"enable_storage_optimizer": false,
			"ignore_file_names":        false,
		},
	}
	_, err := c.Send(data)
	return err
}

func (c *Client) addProxy() {
	p := c.config.proxy
	data := Request{
		"@type":  "addProxy",
		"server": p.host,
		"port":   p.port,
		"enable": true,
		"type": Request{
			"@type":    "proxyTypeSocks5",
			"username": p.login,
			"password": p.password,
		},
	}
	c.SendAndForget(data)
}

func parseResponse(ev Event, req Request, resp interface{}) (err error) {
	err = ev.Unmarshal(resp)
	if err != nil {
		err = ParseErr.Wrap(err, "response parse failed. req: %s", req.String())
	}
	return
}
