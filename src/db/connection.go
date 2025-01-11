package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var Client *mongo.Client
var UserCollection *mongo.Collection

func ConnectDB() {
	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)

	}
	Client = client
	DB = client.Database("CJAssignment")
	UserCollection = DB.Collection("users")
	log.Println("Connected to MongoDB : " + *&clientOptions.Hosts[0])
}

func CloseDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		} else {
			log.Println("MongoDB connection closed.")
		}

	}
}
