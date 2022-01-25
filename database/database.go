package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

func GetError(err error, c *fiber.Ctx) error {
	log.Println(err.Error())

	return GetErrorWithStatus(err, c, fiber.StatusInternalServerError)
}

func GetErrorWithStatus(err error, c *fiber.Ctx, statusCode int) error {
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   statusCode,
	}

	message, _ := json.Marshal(response)

	return c.Status(statusCode).Send(message)
}
