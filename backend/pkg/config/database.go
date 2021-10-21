package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client

func ConnectDatabase() {
	log.Println("Database connecting...")

	//getting DB user info from env variables
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")

	MONGODB_URI := fmt.Sprintf("mongodb+srv://%s:%s@webchat.kcpei.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", user, pass)
	//need to store more parts of the db link in env variables 
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// connecting
	client, err := mongo.Connect(ctx, clientOptions)
	MongoDBClient = client
	if err != nil {
		log.Fatal(err)
	}

	// checking to see if the connection is fine
	err = MongoDBClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected.")
}
