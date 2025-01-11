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

type Database struct {
	Collection *mongo.Collection
}

func NewDatabase() *Database {
	connection = ConnectDB()
	return &Database{GetCollection("character")}
}

var connection *mongo.Client

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx := context.Background()

	// client, err := mongo.NewClient(clientOptions)  // Cria o cliente sem conectar ainda
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = client.Connect(context.Background())  // Conecta explicitamente
	// if err != nil {
	// 	log.Fatal(err)
	// }

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx) // disconnect from the db after function returns

	// checking if the connection succeeded
	err = client.Ping(ctx, readpref.Primary())
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

	// message, _ := json.Marshal(response)

	// return c.Status(statusCode).Send(message)

	return c.Status(statusCode).JSON(response)
}
