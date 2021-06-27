package websocket

import "time"

type StateMessage struct {
	Type       int        `json:"type"`
	ClientList []UserInfo `json:"clientList"`
}

type UserInfo struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	rooms      map[*Room]bool
	// _messageList []Message
	// _messageLimit                 int
	// _expirationLimitHrs           time.Duration
	// _cleanupHeartbeatIntervalMins time.Duration
}

// messageLimit int, expirationLimitHrs time.Duration, cleanupHeartbeatIntervalMins time.Duration

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		rooms:      make(map[*Room]bool),
		// _messageList: []Message{},
		// _messageLimit:                 messageLimit,
		// _expirationLimitHrs:           expirationLimitHrs,
		// _cleanupHeartbeatIntervalMins: cleanupHeartbeatIntervalMins,
	}
}

func (pool *Pool) Start() {
	// go pool.CleanupHeartBeat()
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

	pool.notifyClientJoined(client)
	pool.listClients(client)
	pool.Clients[client] = true

}
func (pool *Pool) unregisterClient(client *Client) {

	if _, ok := pool.Clients[client]; ok {
		delete(pool.Clients, client)
		pool.notifyClientLeft(client)
	}

}

func (pool *Pool) notifyClientJoined(client *Client) {
	message := &Message{
		Action: userJoined,
		// Sender:    client,
		User:      client.User,
		Uid:       client.ID.String(),
		Color:     client.Color,
		Timestamp: time.Now().Format(time.RFC822),
	}
	pool.broadcastToClients(message.encode())
}

func (pool *Pool) notifyClientLeft(client *Client) {
	for room := range client.rooms {
		message := &Message{
			Action: UserLeft,
			// Sender:    client,
			TargetId:  room.ID.String(),
			User:      client.User,
			Uid:       client.ID.String(),
			Timestamp: time.Now().Format(time.RFC822),
		}
		pool.broadcastToClients(message.encode())
	}
}
func (pool *Pool) listClients(client *Client) {
	for existingClient := range pool.Clients {
		message := &Message{
			Action: userJoined,
			// Sender:    existingClient,
			User: existingClient.User,
			// Timestamp: time.Now().Format(time.RFC822),
		}
		client.send <- message.encode()
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
	return foundRoom
}
func (pool *Pool) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	// ^ bool for privacy
	go room.RunRoom()
	pool.rooms[room] = true

	return room
}

func (pool *Pool) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range pool.Clients {
		if client.ID.String() == ID {
			foundClient = client
			break
		}
	}
	return foundClient
}
