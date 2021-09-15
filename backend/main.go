package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/auth"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/config"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/repository"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/websocket"
)

func main() {

	config.ConnectDatabase()
	MongoDb := config.MongoDBClient
	defer MongoDb.Disconnect(context.Background())
	config.CreateRedisClient()

	userRepository := &repository.UserRepository{MongoDB: MongoDb}

	port := os.Getenv("PORT")
	pool := websocket.NewPool(&repository.RoomRepository{MongoDB: MongoDb}, userRepository)
	go pool.Start()

	api := &API{UserRepository: userRepository}

	http.HandleFunc("/ws", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(pool, w, r)
	}))

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)
	if !(port == "") {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+"8080", nil))
	}

	http.HandleFunc("/api/login", api.HandleLogin)
}
