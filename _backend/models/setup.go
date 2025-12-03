package models

import (
	"context" // You need context for Ping
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const db = "wowdle"

var mongoClient *mongo.Client

func ConnectDatabase() {
	connectionString := os.Getenv("MONGO_URI")
	if connectionString == "" {
		connectionString = "mongodb://localhost:27017/"
	}

	clientOption := options.Client().ApplyURI(connectionString)

	// create the client
	client, err := mongo.Connect(clientOption)
	if err != nil {
		panic(err)
	}

	// connect
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// ping database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	mongoClient = client
	println("Successfully connected and pinged MongoDB!")
}
