package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var connection = ConnectDB()

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

func GetCollection(CollectionName string) *mongo.Collection {
	return connection.Database("local").Collection(CollectionName)
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, c *fiber.Ctx) {
	log.Println(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   fiber.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	c.Status(fiber.StatusInternalServerError).Send(message)
}
