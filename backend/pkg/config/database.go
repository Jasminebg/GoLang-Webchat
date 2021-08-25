package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client

func ConnectDatabase() {
	log.Println("Database connecting...")
	// Set client options
	//getting DB user info from .env

	if err := godotenv.Load(); err != nil {
		log.Println("No Env Found")
		// log.Fatal(err)
	}
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")

	MONGODB_URI := fmt.Sprintf("mongodb+srv://%s:%s@webchat.kcpei.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", user, pass)
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
