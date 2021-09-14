package websocket

import (
	"encoding/json"
	"log"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/models"
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
	Sender models.User `json:"sender"`
}

func (message *Message) encode() []byte {
	log.Println("message ")
	log.Println(message)
	if message.Sender != nil {
		log.Println(message.Sender)
		log.Println(message.Sender.GetId())
		log.Println(message.Sender.GetName())

	}

	jsonmessage, err := json.Marshal(message)
	if err != nil {
		log.Println("error ")
		log.Println(err)
	}
	log.Println("jsonmsg ")
	log.Println(jsonmessage)
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
		log.Println("message unmarshal ")
		log.Println(&msg)
		return err
	}
	message.Sender = &msg.Sender
	return nil

}
