package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/config"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/repository"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/websocket"
)

func main() {

	fmt.Println("Chat App ")

	config.ConnectDatabase()
	MongoDb := config.MongoDBClient
	defer MongoDb.Disconnect(context.Background())

	port := os.Getenv("PORT")
	pool := websocket.NewPool(&repository.RoomRepository{MongoDB: MongoDb}, &repository.UserRepository{MongoDB: MongoDb})
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	})

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
