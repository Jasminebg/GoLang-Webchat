package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID    uuid.UUID `json:"id"`
	User  string
	Color string
	Conn  *websocket.Conn
	Pool  *Pool
	send  chan []byte
	rooms map[*Room]bool
}

const (
	maxMessageSize = 1000

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	writeWait = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessageData struct {
	Message string
	Id      string
	Action  string
}

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["user"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url param 'user' is missing")
		return
	}
	color, ok := r.URL.Query()["userColour"]
	if !ok || len(color[0]) < 1 {
		color[0] = "E92750"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(conn, pool, name[0], color[0])

	go client.Write()
	go client.Read()
}

func newClient(conn *websocket.Conn, pool *Pool, name string, color string) *Client {
	return &Client{
		ID:    uuid.New(),
		User:  name,
		Color: color,
		Conn:  conn,
		Pool:  pool,
		send:  make(chan []byte, 256),
		rooms: make(map[*Room]bool),
	}

}

func (client *Client) Read() {
	defer func() {
		client.disconnect()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// func (c *Client) Read() {
// 	defer func() {
// 		c.Pool.Unregister <- c
// 		c.Conn.Close()
// 	}()

// 	for {

// 		messageType, p, err := c.Conn.ReadMessage()

// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		var messageData MessageData
// 		json.Unmarshal([]byte(p), &messageData)

// 		// if messageData.Id != c.ID {
// 		// 	// log.Println(messageData.Id, "./.", c.ID)
// 		// 	log.Println("Unauthorized User")
// 		// 	return
// 		// }

// 		message := Message{
// 			Type:      messageType,
// 			Message:   messageData.Message,
// 			User:      c.User,
// 			Color:     c.Color,
// 			Timestamp: time.Now().Format(time.RFC822)}
// 		// Action:    messageData.Action}

// 		c.handleNewMessage(message.encode())

// 		// c.Pool.Broadcast <- message
// 		fmt.Printf("Message Received: %+v\n", message)
// 	}
// }

func (client *Client) disconnect() {
	client.Pool.Unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	client.Conn.Close()
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}
	message.Sender = client

	switch message.Action {
	case SendMessage:
		roomName := message.Target
		//is room only a local variable here?
		if room := client.Pool.findRoomByName(roomName); room != nil {
			room.broadcast <- &message
		}

	case JoinRoom:
		client.handleJoinRoomMessage(message)

	case LeaveRoom:
		client.handleLeaveRoomMessage(message)

	}

}
func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message
	room := client.Pool.findRoomByName(roomName)
	if room == nil {
		room = client.Pool.createRoom(roomName)
	}
	client.rooms[room] = true
	room.register <- client
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.Pool.findRoomByName(message.Message)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client

}

func (client *Client) GetName() string {
	return client.User
}

// func (client *Client) handleNewMessage(jsonMessage []byte){

// }
