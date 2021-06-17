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
const JoinRoomPrivate = "leave-room"
const RoomJoined = "send-message"

type Message struct {
	Type      int     `json:"type"`
	Message   string  `json:"message"`
	User      string  `json:"user"`
	Color     string  `json:"color"`
	Timestamp string  `json:"timeStamp"`
	Action    string  `json:"action"`
	Target    string  `json:"target"`
	Sender    *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}
