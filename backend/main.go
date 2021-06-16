package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/websocket"
)

type ChatServer struct {
	messageList []websocket.MessageData
}

const (
	Turquoise = "#1ABC9C"
	Orange    = "#E67E2A"
	Red       = "#E92750"
)

func (c *ChatServer) serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	keys := r.URL.Query()
	user := keys.Get("user")
	if len(user) < 1 {
		fmt.Println("Url param 'user' is missing")
		return
	}
	colour := keys.Get("userColour")
	if colour == "" {
		colour = "E92750"
	}
	// else {
	// 	tmp := colour
	// }
	// colour := tmp
	if len(user) < 1 {
		fmt.Println("Url param 'colour' is missing")
		return
	}

	userId := keys.Get("userId")
	if len(userId) < 1 {
		fmt.Println("Url param 'userId' is missing")

	}

	client := &websocket.Client{
		ID:    userId,
		User:  user,
		Color: "#" + colour,
		Conn:  conn,
		Pool:  pool,
	}

	pool.Register <- client
	client.Read()
}
func GetColor() string {
	var colorList = [3]string{Turquoise, Orange, Red}
	rand.Seed(time.Now().Unix())
	return colorList[rand.Intn(3)]
}

func (c *ChatServer) setupRoutes() {
	pool := websocket.NewPool(10, 10, 30)
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c.serveWs(pool, w, r)
	})
}

func main() {
	chatServer := ChatServer{make([]websocket.MessageData, 0)}
	fmt.Println("Chat App v0.01")
	chatServer.setupRoutes()
	http.ListenAndServe(":8080", nil)
}
