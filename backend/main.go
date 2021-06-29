package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/websocket"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	// chatServer := ChatServer{make([]websocket.MessageData, 0)}
	fmt.Println("Chat App v0.1")

	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})
	// chatServer.setupRoutes()
	// use below for deploying?
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
