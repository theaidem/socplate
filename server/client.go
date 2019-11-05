package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	msgpack "github.com/vmihailenco/msgpack/v4"
)

const (
	writeWait      = 1 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	user  *user
	conn  *websocket.Conn
	send  chan *message
	close chan struct{}
	room  *room
}

func (c *client) readPump() {
	defer func() {
		if c.room != nil {
			c.room.unregisterClient <- c
		}
		lobby().unregisterClient <- c
		hub().unregisterClient <- c

		close(c.send)
		close(c.close)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = c.handleMessage(msg)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}

}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok || message == nil {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Fatalln(err)
				return
			}

			msg, err := msgpack.Marshal(message)
			if err != nil {
				log.Fatalln(err)
				return
			}

			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.close:
			return
		}

	}
}

func (c *client) handleMessage(msg []byte) error {

	message, err := parseMessage(msg)
	if err != nil {
		return err
	}

	switch message.Event {
	case test:
		break

	case createRoom:
		if c.room != nil {
			return errors.New("Client already in a room")
		}
		uuid, err := newUUID()
		if err != nil {
			return err
		}
		room := newRoom(uuid)
		room.registerClient <- c
		lobby().registerRoom <- room
		break

	case leaveRoom:
		if c.room == nil {
			return errors.New("Client not in a room")
		}
		c.room.unregisterClient <- c
		if _, ok := lobby().clients.Load(c); !ok {
			lobby().registerClient <- c
		}
		c.send <- newMessage(userEvt, c.user.userResponse())
		break

	case joinRoom:
		if c.room != nil {
			return errors.New("Client already in a room")
		}
		room := lobby().findRoom(message.Payload.(string))
		if room != nil {
			room.registerClient <- c
		} else {
			log.Println("room not found")
		}
		break

	case chatMessage:
		msg := &chatMsg{
			From:    c.user.userResponse(),
			Message: message.Payload.(string),
			Created: time.Now().UnixNano(),
		}
		if _, ok := lobby().clients.Load(c); ok {
			lobby().storeChatMessage(msg)
		} else if c.room != nil {
			c.room.storeChatMessage(msg)
		}
		break

	default:
		return errors.New("Unknown event")
	}

	return nil

}
