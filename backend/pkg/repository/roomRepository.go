package repository

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/models"
)

type Room struct {
	Id      string
	Name    string
	Private bool
}

func (room *Room) GetId() string {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

type RoomRepository struct {
	MongoDBClient *mongo.Client
}

func (repo *RoomRepository) AddRoom(room models.Room) {

}
