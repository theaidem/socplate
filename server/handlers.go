package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {

	ticket := newTicket()
	js, err := json.Marshal(ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func ws(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ticket := query.Get("ticket")
	if !ticketIsValid(ticket) {
		log.Println("Invalid ticket value")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := newUser(rand.Intn((10000 - 1) + 1))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &client{
		user:  user,
		conn:  conn,
		send:  make(chan *message, 256),
		close: make(chan struct{}, 1),
		room:  nil,
	}

	hub().registerClient <- client
	lobby().registerClient <- client

	go client.writePump()
	go client.readPump()
}
