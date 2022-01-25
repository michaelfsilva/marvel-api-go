package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"marvel-api-go/document"
	"marvel-api-go/repository"
)

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

func GetCharacters(c *fiber.Ctx) error {
	log.Println("listing all characters")

	characters := repository.ListAll(c)

	if characters == nil {
		c.SendStatus(fiber.StatusNoContent)
		return nil
	}

	response, _ := json.Marshal(characters) // encode similar to serialize process.
	return c.Send(response)
}

func GetCharacterById(c *fiber.Ctx) error {
	log.Println("listing character by id")

	character := repository.GetById(c)

	if (document.Character{} == character) {
		c.SendStatus(fiber.StatusNoContent)
		return nil
	}

	response, _ := json.Marshal(character)
	return c.Send(response)
}

func GetCharacterByName(c *fiber.Ctx) error {
	log.Println("listing characters by name")

	characters := repository.GetByName(c)

	if characters == nil {
		c.SendStatus(fiber.StatusNoContent)
		return nil
	}

	response, _ := json.Marshal(characters)
	return c.Send(response)
}

func AddCharacter(c *fiber.Ctx) error {
	response, _ := json.Marshal(repository.Add(c))
	return c.Status(fiber.StatusCreated).Send(response)
}

func UpdateCharacter(c *fiber.Ctx) error {
	character := repository.Update(c)

	if (document.Character{} == character) {
		c.SendStatus(fiber.StatusNotFound)
		return nil
	}

	response, _ := json.Marshal(character)
	return c.Send(response)
}

func PartialUpdateCharacter(c *fiber.Ctx) error {
	result := repository.PartialUpdate(c)
	if result == nil {
		return nil
	}

	response, _ := json.Marshal(result)
	return c.Send(response)
}

func DeleteCharacter(c *fiber.Ctx) error {
	jsonResponse, _ := json.Marshal(repository.Delete(c))
	return c.Send(jsonResponse)
}
