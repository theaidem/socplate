package main

import (
	"sync"
)

var (
	lobbyOnce     sync.Once
	lobbyInstance *appLobby
)

type appLobby struct {
	sync.Mutex
	broadcast        chan *message
	clients          sync.Map
	registerClient   chan *client
	unregisterClient chan *client
	rooms            sync.Map
	registerRoom     chan *room
	unregisterRoom   chan *room
	chat             *chat
}

func lobby() *appLobby {
	lobbyOnce.Do(func() {
		lobbyInstance = &appLobby{
			broadcast:        make(chan *message),
			clients:          sync.Map{},
			registerClient:   make(chan *client),
			unregisterClient: make(chan *client),
			rooms:            sync.Map{},
			registerRoom:     make(chan *room),
			unregisterRoom:   make(chan *room),
			chat:             newChat(),
		}
		go lobbyInstance.run()
	})

	return lobbyInstance
}

func (l *appLobby) run() {
	for {
		select {
		case client := <-l.registerClient:
			if _, ok := l.clients.Load(client); !ok {
				l.clients.Store(client, true)
				client.send <- newMessage(roomList, l.availableRooms())
				client.send <- newMessage(chatMessages, l.chat.messageList())
			}

		case client := <-l.unregisterClient:
			if _, ok := l.clients.Load(client); ok {
				l.clients.Delete(client)
				l.broadcastAvailableRooms()
			}

		case room := <-l.registerRoom:
			if _, ok := l.rooms.Load(room); !ok {
				l.rooms.Store(room, true)
				l.broadcastAvailableRooms()
			}

		case room := <-l.unregisterRoom:
			if _, ok := l.rooms.Load(room); ok {
				l.rooms.Delete(room)
				l.broadcastAvailableRooms()
			}

		case message := <-l.broadcast:
			l.clients.Range(func(c, ok interface{}) bool {
				if ok.(bool) {
					c.(*client).send <- message
				} else {
					l.clients.Delete(c)
				}
				return true
			})
		}
	}
}

func (l *appLobby) findRoom(id string) *room {
	var theRoom *room
	l.rooms.Range(func(r, ok interface{}) bool {
		if r.(*room).id == id {
			theRoom = r.(*room)
			return false
		}
		return true
	})
	return theRoom
}

func (l *appLobby) availableRooms() []*roomResponse {
	availableRooms := []*roomResponse{}
	l.rooms.Range(func(r, ok interface{}) bool {
		// Clean up empty rooms
		if r.(*room).clientsCount() == 0 {
			l.rooms.Delete(r)
			return true
		}
		if r.(*room).clientsCount() < maxClients && !r.(*room).isStarted {
			users := []*userResponse{}
			r.(*room).clients.Range(func(c, ok interface{}) bool {
				user := c.(*client).user.userResponse().addRoom(r.(*room).id)
				users = append(users, user)
				return true
			})
			room := &roomResponse{
				ID:    r.(*room).id,
				Users: users,
			}
			availableRooms = append(availableRooms, room)
		}
		return true
	})
	return availableRooms
}

func (l *appLobby) broadcastAvailableRooms() {
	go func() { l.broadcast <- newMessage(roomList, l.availableRooms()) }()
}

func (l *appLobby) roomsCount() int64 {
	var length int64
	l.rooms.Range(func(key interface{}, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (l *appLobby) broadcastRoomsCount() {
	go func() { l.broadcast <- newMessage(roomsCount, l.roomsCount()) }()
}

func (l *appLobby) storeChatMessage(message *chatMsg) {
	l.chat.storeChatMessage(message)
	l.broadcastChatMessages()
}

func (l *appLobby) broadcastChatMessages() {
	go func() { l.broadcast <- newMessage(chatMessages, l.chat.messageList()) }()
}
