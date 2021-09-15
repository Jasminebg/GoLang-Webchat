package repository

import (
	"context"
	"os"
	"time"

	"github.com/Jasminebg/GoLang-Webchat/backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id       string `json:"id"`
	User     string `json:"user"`
	Password string `json:"password"`
	Color    string `json:"color"`
}

func (user *User) GetName() string {
	return user.User
}
func (user *User) GetId() string {
	return user.Id
}

func (user *User) GetColor() string {
	return user.Color
}

type UserRepository struct {
	MongoDB *mongo.Client
}

func (repo *UserRepository) AddUser(user models.User) {
	collection := repo.MongoDB.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, registrationError := collection.InsertOne(ctx, bson.M{
		"userId":   user.GetId(),
		"userName": user.GetName(),
	})

	defer cancel()

	checkErr(registrationError)
}

func (repo *UserRepository) RemoveUser(user models.User) {
	collection := repo.MongoDB.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_, registrationError := collection.DeleteOne(ctx, bson.M{
		"userId":   user.GetId(),
		"userName": user.GetName(),
	})

	defer cancel()

	checkErr(registrationError)
}
func (repo *UserRepository) FindUserByUsername(userName string) *User {
	collection := repo.MongoDB.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user User
	//find room by name
	err := collection.FindOne(ctx, bson.M{"userName": userName}).Decode(&user)

	checkErr(err)

	//assign room id, name and private to returned struct, check for err
	defer cancel()

	return &user
}

func (repo *UserRepository) FindUserById(ID string) models.User {
	collection := repo.MongoDB.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user User
	//find room by name
	err := collection.FindOne(ctx, bson.M{"userId": ID}).Decode(&user)

	checkErr(err)

	//assign room id, name and private to returned struct, check for err
	defer cancel()

	return &user
}

func (repo *UserRepository) GetAllUsers() []models.User {

	var users []models.User

	collection := repo.MongoDB.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})

	defer cancel()

	if err != nil {
		return users
	}

	for cursor.Next(context.TODO()) {
		var user models.User
		err := cursor.Decode(&user)
		if err == nil {
			users = append(users, user)
		}
	}

	return users
}
