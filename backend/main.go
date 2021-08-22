package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/config"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/websocket"
)

func main() {

	fmt.Println("Chat App ")

	port := os.Getenv("PORT")
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})

	config.ConnectDatabase()
	MongoDB := config.MongoDBClient
	defer MongoDB.Close()

	// use below for deploying?

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)
	if !(port == "") {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+"8080", nil))
	}
	// log.Fatal(http.ListenAndServe(*addr, nil))
}
