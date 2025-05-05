package main

import (
	"log"
	"marvel-api-go/app"
	"marvel-api-go/repository"
)

func main() {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository("mongodb://localhost:27017")

	app := app.SetupApp(&repository)

	log.Fatal(app.Listen(":8080"))
}
