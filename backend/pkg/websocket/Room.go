package websocket

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const welcomeMessage = "%s joined the room"

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	Private    bool `json:"private"`
}

func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		Private:    private,
	}

}

func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			message.Timestamp = time.Now().Format(time.RFC822)
			room.broadcastToClientsInRoom(message.encode())

		}

	}

}

func (room *Room) registerClientInRoom(client *Client) {
	if !room.Private {
		room.notifyClientJoined(client)
		room.clients[client] = true
		room.listClientsinRoom(client)

	}

}
func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}

}
func (room *Room) broadcastToClientsInRoom(message []byte) {
	fmt.Println(room.clients)
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) listClientsinRoom(client *Client) {
	for otherclient := range room.clients {
		if otherclient.GetID() != client.GetID() {
			message := &Message{
				Action:   userJoinedRoom,
				User:     otherclient.GetName(),
				Color:    otherclient.GetColor(),
				Uid:      otherclient.GetID(),
				TargetId: room.ID.String(),
			}
			client.send <- message.encode()

		}
	}
}
func (room *Room) notifyClientJoined(sender *Client) {
	message := &Message{
		Message:   fmt.Sprintf(welcomeMessage, sender.GetName()),
		Action:    SendMessage,
		Target:    room.Name,
		TargetId:  room.ID.String(),
		Timestamp: time.Now().Format(time.RFC822),
	}
	room.broadcastToClientsInRoom(message.encode())

	joinMessage := &Message{
		Action:   userJoinedRoom,
		User:     sender.GetName(),
		Color:    sender.GetColor(),
		Uid:      sender.GetID(),
		Target:   room.Name,
		TargetId: room.ID.String(),
	}
	room.broadcastToClientsInRoom(joinMessage.encode())

}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}
