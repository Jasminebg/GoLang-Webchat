package websocket

import (
	"encoding/json"
	"log"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/config"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/models"
	"github.com/google/uuid"
)

type StateMessage struct {
	Type       int        `json:"type"`
	ClientList []UserInfo `json:"clientList"`
}

type UserInfo struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Pool struct {
	Register       chan *Client
	Unregister     chan *Client
	Clients        map[*Client]bool
	Broadcast      chan Message
	rooms          map[*Room]bool
	users          []models.User
	roomRepository models.RoomRepository
	userRepository models.UserRepository
	// _messageList []Message
	// _messageLimit                 int
	// _expirationLimitHrs           time.Duration
	// _cleanupHeartbeatIntervalMins time.Duration
}

const PubSubGeneralChannel = "general"

// messageLimit int, expirationLimitHrs time.Duration, cleanupHeartbeatIntervalMins time.Duration

func NewPool(roomRepository models.RoomRepository, userRepository models.UserRepository) *Pool {
	pool := &Pool{
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Clients:        make(map[*Client]bool),
		Broadcast:      make(chan Message),
		rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
		// _messageList: []Message{},
		// _messageLimit:                 messageLimit,
		// _expirationLimitHrs:           expirationLimitHrs,
		// _cleanupHeartbeatIntervalMins: cleanupHeartbeatIntervalMins,
	}
	pool.users = userRepository.GetAllUsers()

	return pool
}

func (pool *Pool) Start() {
	// go pool.CleanupHeartBeat()
	go pool.listenPubSubChannel()
	for {
		select {
		//connecting
		case client := <-pool.Register:
			pool.registerClient(client)

			break
		//disconnecting
		case client := <-pool.Unregister:
			pool.unregisterClient(client)
			break
		//broadcasting message
		case message := <-pool.Broadcast:
			pool.broadcastToClients(message.encode())

		}
	}
}

func (pool *Pool) registerClient(client *Client) {

	if user := pool.findUserByID(client.ID); user == nil {
		pool.userRepository.AddUser(client)
	}
	pool.publishClientJoined(client)
	pool.listClients(client)
	pool.Clients[client] = true

	pool.users = append(pool.users, client)

}
func (pool *Pool) unregisterClient(client *Client) {

	if _, ok := pool.Clients[client]; ok {
		delete(pool.Clients, client)
		// pool.userRepository.RemoveUser(client)
		pool.publishClientLeft(client)

	}

}

func (pool *Pool) publishClientJoined(client *Client) {
	message := &Message{
		Action: userJoined,
		User:   client.User,
		Uid:    client.ID,
		// Sender: client,
	}
	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (pool *Pool) publishClientLeft(client *Client) {
	message := &Message{
		Action: UserLeft,
		User:   client.User,
		Uid:    client.ID,
		// Sender: client,
	}
	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (pool *Pool) listenPubSubChannel() {
	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)
	ch := pubsub.Channel()
	for msg := range ch {
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Err on unmarshal %s", err)
			return
		}
		switch message.Action {
		case userJoined:
			pool.handleUserJoined(message)
		case UserLeft:
			pool.handleUserLeft(message)
		case JoinRoomPrivate:
			pool.handleJoinRoomPrivateMessage(message)
		}
	}
}

func (pool *Pool) handleUserJoined(message Message) {
	pool.users = append(pool.users, pool.findClientByID(message.Uid))
	pool.broadcastToClients(message.encode())
}

func (pool *Pool) handleUserLeft(message Message) {
	for i, user := range pool.users {
		if user.GetId() == message.Uid {
			pool.users[i] = pool.users[len(pool.users)-1]
			pool.users = pool.users[:len(pool.users)-1]
			break
		}
	}
	pool.broadcastToClients(message.encode())
}

func (pool *Pool) handleJoinRoomPrivateMessage(message Message) {
	targetClients := pool.findClientsByID(message.Message)
	for _, targetClient := range targetClients {
		targetClient.joinRoom(message.Target, pool.findClientByID(message.Uid))
	}

}

func (pool *Pool) findClientsByID(ID string) []*Client {
	var foundClients []*Client
	for client := range pool.Clients {
		if client.GetId() == ID {
			foundClients = append(foundClients, client)
		}
	}
	return foundClients
}

func (pool *Pool) findUserByID(ID string) models.User {
	var foundUser models.User
	for _, client := range pool.users {
		if client.GetId() == ID {
			foundUser = client
			break
		}
	}
	return foundUser
}

func (pool *Pool) listClients(client *Client) {
	var uniqueUsers = make(map[string]bool)
	for _, user := range pool.users {
		if ok := uniqueUsers[user.GetId()]; !ok {
			message := &Message{
				Action: userJoined,
				User:   user.GetName(),
				Uid:    user.GetId(),
				Color:  user.GetColor(),
			}
			uniqueUsers[user.GetId()] = true
			client.send <- message.encode()
		}
	}
}

func (pool *Pool) broadcastToClients(message []byte) {

	for client := range pool.Clients {
		client.send <- message
	}
}

func (pool *Pool) GetClientNames() []UserInfo {
	clients := make([]UserInfo, len(pool.Clients))
	i := 0
	for k := range pool.Clients {
		clients[i] = UserInfo{
			Name:  k.User,
			Color: k.Color,
		}
		i++
	}
	return clients
}

func (pool *Pool) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range pool.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

func (pool *Pool) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range pool.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	if foundRoom == nil {
		foundRoom = pool.runRoomFromRepository(name)
	}
	if foundRoom.Name == "" {
		return nil
	}

	return foundRoom
}

func (pool *Pool) runRoomFromRepository(name string) *Room {
	var room *Room

	dbRoom := pool.roomRepository.FindRoomByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.RunRoom()
		pool.rooms[room] = true
	}

	return room
}

func (pool *Pool) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	// ^ bool for privacy
	pool.roomRepository.AddRoom(room)
	go room.RunRoom()
	pool.rooms[room] = true

	return room
}

func (pool *Pool) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range pool.Clients {
		if client.ID == ID {
			foundClient = client
			break
		}
	}
	return foundClient
}
