package websocket

import (
	"encoding/json"
	"log"
)

const SendMessage = "send-message"
const JoinRoom = "join-room"
const LeaveRoom = "leave-room"
const userJoined = "user-join"
const userJoinedRoom = "user-join-room"
const listRoomClients = "list-clients"
const UserLeft = "user-left"
const JoinRoomPrivate = "join-room-private"
const RoomJoined = "room-joined"

type Message struct {
	// Type      int     `json:"type"`
	Message string `json:"message"`
	// Sender    models.User `json:"sender"`
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	// Room      *Room       `json:"target"`
	Private  bool   `json:"private"`
	User     string `json:"user"`
	Uid      string `json:"id"`
	Color    string `json:"color"`
	Target   string `json:"room"`
	TargetId string `json:"roomid"`
}

func (message *Message) encode() []byte {
	// log.Println("message ")
	// log.Println(message)

	jsonmessage, err := json.Marshal(message)
	if err != nil {
		log.Println("error ")
		log.Println(err)
	}
	return jsonmessage
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		// log.Println("message unmarshal ")
		// log.Println(&msg)
		return err
	}
	// message.Sender = &msg.Sender
	return nil

}
