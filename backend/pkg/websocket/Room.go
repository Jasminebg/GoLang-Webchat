package websocket

import (
	"fmt"

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
			room.broadcastToClientsInRoom(message.encode())

		}

	}

}

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true

}
func (room *Room) unregisterClientInRoom(client *Client) {

}
func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}

}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessage,
		Target:  room.Name,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}
	room.broadcastToClientsInRoom(message.encode())

}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}
