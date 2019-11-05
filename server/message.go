package main

import msgpack "github.com/vmihailenco/msgpack/v4"

type event string

const (
	test         event = "test"
	err          event = "err"
	warn         event = "warn"
	online       event = "online"
	userEvt      event = "user"
	roomsCount   event = "rooms_count"
	roomList     event = "room_list"
	chatMessage  event = "chat_message"
	chatMessages event = "chat_messages"
	createRoom   event = "create_room"
	leaveRoom    event = "leave_room"
	joinRoom     event = "join_room"
	startRoom    event = "start_room"
	updateRoom   event = "update_room"
)

type message struct {
	Event   event       `msgpack:"event"`
	Payload interface{} `msgpack:"payload"`
}

func newMessage(evt event, payload interface{}) *message {
	return &message{
		Event:   evt,
		Payload: payload,
	}
}

func parseMessage(msg []byte) (*message, error) {
	message := message{}
	if err := msgpack.Unmarshal(msg, &message); err != nil {
		return nil, err
	}
	return &message, nil
}
