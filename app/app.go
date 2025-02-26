package app

import (
	"marvel-api-go/controller"
	"marvel-api-go/repository"
	"marvel-api-go/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func SetupApp(repository *repository.CharacterRepositoryImpl) *fiber.App {
	app := fiber.New()

	// Provide a minimal auth config
	// https://docs.gofiber.io/api/middleware/basicauth
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"test": "marvel",
		},
	}))

	service := service.NewCharacterService(repository)
	controller := controller.NewCharacterController(*service)

	setupRoutes(app, controller)

	return app
}

func setupRoutes(app *fiber.App, controller *controller.CharacterController) {
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-type", "application/json")
		return c.Next()
	})

	app.Get("/api/characters", controller.GetAllCharacters)
	//app.Get("/api/characters/:id?", controller.GetAllCharactersOrFilterById)
	app.Get("/api/characters/:id", controller.GetCharacterById)
	app.Get("/api/characters/findByName/:name", controller.GetCharacterByName)
	app.Post("/api/characters", controller.AddCharacter)
	app.Put("/api/characters/:id", controller.UpdateCharacter)
	app.Patch("/api/characters/:id", controller.PartialUpdateCharacter)
	app.Delete("/api/characters/:id", controller.DeleteCharacter)
}
