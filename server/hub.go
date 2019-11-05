package main

import (
	"sync"
)

var (
	hubOnce     sync.Once
	hubInstance *appHub
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type appHub struct {
	clients          sync.Map
	broadcast        chan *message
	registerClient   chan *client
	unregisterClient chan *client
}

func hub() *appHub {
	hubOnce.Do(func() {
		hubInstance = &appHub{
			clients:          sync.Map{},
			broadcast:        make(chan *message),
			registerClient:   make(chan *client),
			unregisterClient: make(chan *client),
		}
		go hubInstance.run()
	})

	return hubInstance
}

func (h *appHub) run() {
	for {
		select {
		case client := <-h.registerClient:
			h.clients.Store(client, true)
			client.send <- newMessage(userEvt, client.user.userResponse())
			h.broadcastOnlineCount()

		case client := <-h.unregisterClient:
			if _, ok := h.clients.Load(client); ok {
				h.clients.Delete(client)
				h.broadcastOnlineCount()
			}

		case message := <-h.broadcast:
			h.clients.Range(func(c, ok interface{}) bool {
				if ok.(bool) {
					c.(*client).send <- message
				} else {
					h.clients.Delete(c)
					close(c.(*client).send)
				}
				return true
			})

		}
	}
}

func (h *appHub) close() {
	h.clients.Range(func(c, ok interface{}) bool {
		if ok.(bool) {
			h.unregisterClient <- c.(*client)
		}
		return true
	})
}

func (h *appHub) clientsCount() int64 {
	var length int64
	h.clients.Range(func(key interface{}, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (h *appHub) broadcastOnlineCount() {
	go func() { h.broadcast <- newMessage(online, h.clientsCount()) }()
}
