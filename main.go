package main

import (
	"github.com/gofiber/fiber"
	"marvel-api-go/controller"
)

func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) {
		c.Set("Content-type", "application/json")
		c.Next()
	})

	app.Get("/api/characters", controller.GetCharacters)
	//app.Get("/api/characters/:id?", controller.GetAllCharactersOrFilterById)
	app.Get("/api/characters/:id", controller.GetCharacterById)
	app.Get("/api/characters/findByName/:name", controller.GetCharacterByName)
	app.Post("/api/characters", controller.AddCharacter)
	app.Put("/api/characters/:id", controller.UpdateCharacter)
	app.Patch("/api/characters/:id", controller.PartialUpdateCharacter)
	app.Delete("/api/characters/:id", controller.DeleteCharacter)

	app.Listen("8080")
}
