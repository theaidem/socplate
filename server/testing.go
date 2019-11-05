package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const clientsCount = 1500

func generateClients() {
	for index := 0; index < clientsCount; index++ {

		// Create test server with the echo handler.
		s := httptest.NewServer(http.HandlerFunc(ws))
		defer s.Close()

		ticket := newTicket()

		u := "ws" + strings.TrimPrefix(s.URL, "http")
		u = fmt.Sprintf("%s?ticket=%s", u, ticket.Tiket)

		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			log.Printf("%v", err)
		}
		defer ws.Close()

		<-time.Tick(time.Millisecond * 500)

	}
	<-time.Tick(time.Second * 120)
}
