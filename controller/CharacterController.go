package controller

import (
	"encoding/json"
	"log"
	"marvel-api-go/document"
	. "marvel-api-go/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharacterController struct {
	characterService CharacterService
}

func NewCharacterController(service CharacterService) *CharacterController {
	return &CharacterController{characterService: service}
}

//func GetAllCharactersOrFilterById(c *fiber.Ctx) error {
//	var characters []document.Character
//	var filter = bson.M{}
//
//	// this if makes this methods dynamic to get all if the id is not found
//	if c.Params("id") != "" {
//		id := c.Params("id")
//		objID, _ := primitive.ObjectIDFromHex(id)
//		filter = bson.M{"_id": objID}
//	}
//
//	cur, err := collection.Find(context.Background(), filter)
//	defer cur.Close(context.Background())
//
//	if err != nil {
//		database.GetError(err, c)
//		return
//	}
//
//	for cur.Next(context.Background()) {
//		var character document.Character
//
//		// & returns the memory address of the following variable.
//		err := cur.Decode(&character) // decode similar to deserialize process.
//		if err != nil {
//			database.GetError(err, c)
//		}
//
//		characters = append(characters, character)
//	}
//
//	if err := cur.Err(); err != nil {
//		database.GetError(err, c)
//	}
//
//	response, _ := json.Marshal(characters)
//	c.Send(response)
//}

func (c *CharacterController) GetAllCharacters(ctx *fiber.Ctx) error {
	log.Println("listing all characters")

	characters, err := c.characterService.ListAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if characters == nil {
		return ctx.SendStatus(fiber.StatusNoContent)
	}

	response, _ := json.Marshal(characters) // encode similar to serialize process.
	return ctx.Send(response)
}

func (c *CharacterController) GetCharacterById(ctx *fiber.Ctx) error {
	log.Println("listing character by id")

	character, err := c.characterService.GetCharacterById(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if (document.Character{} == *character) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Character not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(character)
}

func (c *CharacterController) GetCharacterByName(ctx *fiber.Ctx) error {
	log.Println("listing characters by name")

	characters, err := c.characterService.GetCharacterByName(ctx.Params("name"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if characters == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(characters)
}

func (c *CharacterController) AddCharacter(ctx *fiber.Ctx) error {
	var character document.Character
	json.Unmarshal(ctx.Body(), &character)

	serviceResponse, _ := c.characterService.AddCharacter(character)
	response, _ := json.Marshal(serviceResponse)

	return ctx.Status(fiber.StatusCreated).Send(response)
}

func (c *CharacterController) UpdateCharacter(ctx *fiber.Ctx) error {
	var character document.Character
	json.Unmarshal(ctx.Body(), &character)
	character.ID, _ = primitive.ObjectIDFromHex(ctx.Params("id"))

	serviceResponse, err := c.characterService.UpdateCharacter(character)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if (document.Character{} == serviceResponse) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	response, _ := json.Marshal(serviceResponse)
	return ctx.Send(response)
}

func (c *CharacterController) PartialUpdateCharacter(ctx *fiber.Ctx) error {
	var character document.Character
	json.Unmarshal(ctx.Body(), &character)

	serviceResponse, err := c.characterService.PartialUpdateCharacter(character)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if (document.Character{} == serviceResponse) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	response, _ := json.Marshal(serviceResponse)
	return ctx.Send(response)
}

func (c *CharacterController) DeleteCharacter(ctx *fiber.Ctx) error {
	jsonResponse, err := json.Marshal(c.characterService.DeleteCharacter(ctx.Params("id")))
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Send(jsonResponse)
}
