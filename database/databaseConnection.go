package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initiateConnection() *mongo.Client {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var uri = os.Getenv("URI")
	fmt.Println(uri)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err != client.Ping(ctx, nil) {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = initiateConnection()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("socmed").Collection(collectionName)

	return collection
}
