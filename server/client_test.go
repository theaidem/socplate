package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestExample(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(ws))
	defer s.Close()

	ticket := newTicket()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	u = fmt.Sprintf("%s?ticket=%s", u, ticket.Tiket)

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	<-time.Tick(time.Second * 1)

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg, err := parseMessage(p)
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Println(msg)
	if msg.Event != userEvt {
		t.Fatalf("bad message")
	}

	_, p, err = ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg, err = parseMessage(p)
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Println(msg)
	if msg.Event != online {
		t.Fatalf("bad message")
	}

}
