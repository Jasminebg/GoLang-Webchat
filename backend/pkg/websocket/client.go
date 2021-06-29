package websocket

import (
	"encoding/json"
	"fmt"
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
		var ms Message
		if err := json.Unmarshal(jsonMessage, &ms); err != nil {
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
			// message.Timestamp = time.Now().Format(time.RFC822)
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
			var ms Message
			if err := json.Unmarshal(message, &ms); err != nil {
			}
			fmt.Println("unmarshal")
			fmt.Println(ms)
			// fmt.Println(ms.Message)

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
	client.Pool.Register <- client
}
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
		return
	}
	message.User = client.User
	message.Uid = client.ID.String()
	message.Color = client.Color
	message.Timestamp = time.Now().Format(time.RFC822)

	switch message.Action {
	case SendMessage:
		// roomID := client.Pool.findRoomByName()
		//is room only a local variable here?
		// message.Color = client.Color
		if room := client.Pool.findRoomByID(message.TargetId); room != nil {
			room.broadcast <- &message
			// fmt.Println("send message")
			// fmt.Println(room, message)
		}

	case JoinRoom:
		client.handleJoinRoomMessage(message)

	case LeaveRoom:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivate:
		client.handleJoinRoomPrivateMessage(message)
	}
}
func (client *Client) handleJoinRoomMessage(message Message) {
	client.joinRoom(message.Message, nil)
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.Pool.findRoomByName(message.Message)
	if room == nil {
		return
	}

	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client

}

func (client *Client) handleJoinRoomPrivateMessage(message Message) {
	target := client.Pool.findClientByID(message.Message)

	if target == nil {
		return
	}
	roomName := target.User + " and " + client.User + " PMs"

	client.joinRoom(roomName, target)
	target.joinRoom(roomName, client)

}
func (client *Client) joinRoom(roomName string, sender *Client) {
	room := client.Pool.findRoomByName(roomName)
	if room == nil {
		room = client.Pool.createRoom(roomName, sender != nil)
	}

	if sender == nil && room.Private {
		return
	}
	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
		client.notifyRoomJoined(room)
		// if sender == nil {
		// } else {
		// 	client.notifyPrivateRoomJoined(room, sender)
		// }

	}
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false

}

func (client *Client) notifyRoomJoined(room *Room) {
	message := Message{
		Action:   RoomJoined,
		Target:   room.Name,
		TargetId: room.ID.String(),
		// Sender: client,
		User:  client.User,
		Color: client.Color,
		Uid:   client.ID.String(),
	}
	client.send <- message.encode()

}
func (client *Client) notifyPrivateRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action:   RoomJoined,
		Target:   room.Name,
		TargetId: room.ID.String(),
		// Sender: client,
		User:  client.User,
		Color: client.Color,
		Uid:   client.ID.String(),
	}
	client.send <- message.encode()

	// userMessage := &Message{
	// 	Action:   userJoinedRoom,
	// 	User:     sender.GetName(),
	// 	Color:    sender.GetColor(),
	// 	Uid:      sender.GetID(),
	// 	TargetId: room.ID.String(),
	// }
	// client.send <- userMessage.encode()

}

func (client *Client) GetName() string {
	return client.User
}
func (client *Client) GetColor() string {
	return client.Color
}

func (client *Client) GetID() string {
	return client.ID.String()
}

// func (client *Client) handleNewMessage(jsonMessage []byte){

// }
// func (client *Client) handleNewMessage(jsonMessage []byte){

// }
