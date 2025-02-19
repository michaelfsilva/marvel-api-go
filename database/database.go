package database

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
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
	// TODO check the line below
	// defer client.Disconnect(context.Background()) // disconnect from the db after function returns

	// checking if the connection succeeded
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database("local").Collection(collectionName)
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

	// message, _ := json.Marshal(response)
	// return c.Status(statusCode).Send(message)

	return c.Status(statusCode).JSON(response) // this does the same as above
}
