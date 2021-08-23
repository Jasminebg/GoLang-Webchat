package repository

import (
	"context"
	"os"
	"time"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/config"
	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
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

// type RoomRepository struct {
// 	MongoDB := config.MongoDBClient
// }

func AddRoom(room models.Room) {
	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, registrationError := collection.InsertOne(ctx, bson.M{
		"roomID":      room.GetId(),
		"roomName":    room.GetName(),
		"roomPrivate": room.GetPrivate(),
	})

	defer cancel()

	checkErr(registrationError)
}

func FindRoomByName(name string) models.Room {
	collection := config.MongoDBClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("rooms")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var room Room
	//find room by name
	err := collection.FindOne(ctx, bson.M{"roomName": name}).Decode(&room)

	checkErr(err)

	//assign room id, name and private to returned struct, check for err
	defer cancel()

	return &room
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
