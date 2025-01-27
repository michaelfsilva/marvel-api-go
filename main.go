package main

import (
	"log"
	"marvel-api-go/controller"
	"marvel-api-go/repository"
	"marvel-api-go/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func main() {
	app := fiber.New()

	// Provide a minimal config
	// https://docs.gofiber.io/api/middleware/basicauth
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"test": "marvel",
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-type", "application/json")
		return c.Next()
	})

	repository := &repository.CharacterRepositoryImpl{}
	repository.InitRepository()
	service := service.NewCharacterService(repository)
	controller := controller.NewCharacterController(*service)

	app.Get("/api/characters", controller.GetCharacters)
	//app.Get("/api/characters/:id?", controller.GetAllCharactersOrFilterById)
	app.Get("/api/characters/:id", controller.GetCharacterById)
	app.Get("/api/characters/findByName/:name", controller.GetCharacterByName)
	app.Post("/api/characters", controller.AddCharacter)
	app.Put("/api/characters/:id", controller.UpdateCharacter)
	app.Patch("/api/characters/:id", controller.PartialUpdateCharacter)
	app.Delete("/api/characters/:id", controller.DeleteCharacter)

	log.Fatal(app.Listen(":8080"))
}
