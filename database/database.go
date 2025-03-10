package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Collection *mongo.Collection

func NewDatabase(connectionString string, collectionName string) {
	Collection = ConnectDB(connectionString, collectionName)
}

func ConnectDB(connectionString string, collectionName string) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(connectionString)

	// client, err := mongo.NewClient(clientOptions)  // creates the client without connecting yet
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = client.Connect(context.Background())  // connect to the database
	// if err != nil {
	// 	log.Fatal(err)
	// }

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// this does not work using dependency injection
	// defer client.Disconnect(context.Background()) // disconnect from the db after function returns

	// checking if the connection succeeded
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database("local").Collection(collectionName)
}
