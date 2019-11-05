package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

const port = ":3003"

func main() {

	log.Println(database())

	setrLimits()
	http.HandleFunc("/index", index)
	http.HandleFunc("/ws", ws)

	go func() {
		log.Println("Server start, port:", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	hub().close()

}
