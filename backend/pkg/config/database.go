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
	// MONGODB_URI := os.Getenv("MONGODB_URI")
	// if MONGODB_URI == "" {
	// 	MONGODB_URI = "mongodb://localhost/chat"
	// }
	//getting DB user info from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
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

// func initDB(){
// 	MONGODB_URI := os.Getenv("MONGODB_URI")
// 	if MONGODB_URI == ""{
// 		MONGODB_URI = "mongodb://localhost/chat"
// 	}
// 	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_URI))
// 	if err != nil{
// 		log.Fatal(err)
// 	}

// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	err - client.Connect(ctx)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)

// 	databases, err := client.ListDatabaseNames(ctx, bson.M{})
// 	if err != nil {
// 			log.Fatal(err)
// 	}
// 	fmt.Println(databases)

// }
