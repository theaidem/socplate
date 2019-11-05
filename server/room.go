package main

import "sync"

const maxClients = 2

type room struct {
	id               string
	broadcast        chan *message
	clients          sync.Map
	registerClient   chan *client
	unregisterClient chan *client
	chat             *chat
	isStarted        bool
}

type roomResponse struct {
	ID    string          `msgpack:"id"`
	Users []*userResponse `msgpack:"users"`
}

func newRoom(id string) *room {
	room := &room{
		id:               id,
		broadcast:        make(chan *message),
		clients:          sync.Map{},
		registerClient:   make(chan *client),
		unregisterClient: make(chan *client),
		chat:             newChat(),
		isStarted:        false,
	}
	go room.run()
	return room
}

func (r *room) run() {
	for {
		select {
		case client := <-r.registerClient:
			if _, ok := r.clients.Load(client); !ok {
				client.room = r
				r.clients.Store(client, true)
				client.send <- newMessage(userEvt, client.user.userResponse().addRoom(r.id))
				if !r.isStarted {
					lobby().broadcastAvailableRooms()
				}
				if r.clientsCount() == maxClients {
					go r.startRoom()
				}
			}

		case client := <-r.unregisterClient:
			if _, ok := r.clients.Load(client); ok {
				r.clients.Delete(client)
				client.room = nil
				r.broadcastRoomUpdates()
				if !r.isStarted {
					lobby().broadcastAvailableRooms()
				}
			}

		case message := <-r.broadcast:
			r.clients.Range(func(c, ok interface{}) bool {
				if ok.(bool) {
					c.(*client).send <- message
				} else {
					r.clients.Delete(c)
				}
				return true
			})
		}
	}
}

func (r *room) clientsCount() int64 {
	var length int64
	r.clients.Range(func(key interface{}, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (r *room) broadcastRoomUpdates() {
	if r.isStarted {
		go func() { r.broadcast <- newMessage(updateRoom, r.roomResponse()) }()
	}
}

func (r *room) startRoom() {
	r.isStarted = true
	r.clients.Range(func(c, ok interface{}) bool {
		lobby().unregisterClient <- c.(*client)
		return true
	})
	r.broadcast <- newMessage(startRoom, r.roomResponse())
}

func (r *room) roomResponse() *roomResponse {
	users := []*userResponse{}
	r.clients.Range(func(c, ok interface{}) bool {
		users = append(users, c.(*client).user.userResponse().addRoom(r.id))
		return true
	})
	return &roomResponse{
		ID:    r.id,
		Users: users,
	}
}

func (r *room) storeChatMessage(message *chatMsg) {
	r.chat.storeChatMessage(message)
	r.broadcastChatMessages()
}

func (r *room) broadcastChatMessages() {
	go func() { r.broadcast <- newMessage(chatMessages, r.chat.messageList()) }()
}
