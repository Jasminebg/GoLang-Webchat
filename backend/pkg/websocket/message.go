package websocket

import (
	"encoding/json"
	"log"
)

const SendMessage = "send-message"
const JoinRoom = "join-room"
const LeaveRoom = "leave-room"
const userJoined = "send-message"
const UserLeft = "user-left"
const JoinRoomPrivate = "join-room-private"
const RoomJoined = "room-joined"

type Message struct {
	// Type      int     `json:"type"`
	Message   string `json:"message"`
	User      string `json:"user"`
	Uid       string `json:"id"`
	Color     string `json:"color"`
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Target    string `json:"room"`
	TargetId  string `json:"roomid"`
	Private   bool   `json:"private"`
	// Target    *Room   `json:"target"`
	// Sender    *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}
