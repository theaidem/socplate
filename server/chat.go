package main

import (
	"sync"
)

type chat struct {
	sync.Mutex
	messages []*chatMsg
}

type chatMsg struct {
	From    *userResponse `msgpack:"from"`
	Message string        `msgpack:"message"`
	Created int64         `msgpack:"created"`
}

func newChat() *chat {
	return &chat{
		messages: []*chatMsg{},
	}
}

func (c *chat) storeChatMessage(message *chatMsg) {
	c.Lock()
	defer c.Unlock()
	c.messages = append(c.messages, message)
	return
}

func (c *chat) messageList() []*chatMsg {
	c.Lock()
	defer c.Unlock()
	return c.messages
}
