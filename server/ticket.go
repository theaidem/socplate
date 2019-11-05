package main

import (
	"sync"
	"time"
)

var (
	tickets   = sync.Map{}
	ttlTicket = 5 * time.Second
)

type ticket struct {
	Tiket   string `json:"ticket"`
	Exp     int64  `json:"exp"`
	Created int64  `json:"created"`
}

func newTicket() *ticket {
	t := &ticket{
		Tiket:   randomString(10),
		Exp:     ttlTicket.Milliseconds(),
		Created: time.Now().Unix(),
	}

	tickets.Store(t, true) //[t] = true
	go func(t *ticket) {
		time.Sleep(ttlTicket)
		tickets.Delete(t)
	}(t)

	return t
}

func ticketIsValid(t string) bool {
	found := false
	tickets.Range(func(key interface{}, value interface{}) bool {
		if key.(*ticket).Tiket == t {
			tickets.Delete(key)
			found = true
			return false
		}
		return true
	})
	return found
}
