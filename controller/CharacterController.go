package controller

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"marvel-api-go/database"
	"marvel-api-go/document"
)

//Connection mongoDB
var collection = database.GetCollection("character")

//func GetAllCharactersOrFilterById(c *fiber.Ctx) {
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

func GetCharacters(c *fiber.Ctx) {
	var characters []document.Character

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.Background(), bson.M{})

	// Close the cursor once finished
	// A defer statement defers the execution of a function until the surrounding function returns.
	// simply, run cur.Close() process but after cur.Next() finished.
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.Background())

	if err != nil {
		database.GetError(err, c)
		return
	}

	// better than using a loop
	err = cur.All(context.Background(), &characters)

	if err != nil || characters == nil {
		c.SendStatus(404)
		return
	}

	response, _ := json.Marshal(characters) // encode similar to serialize process.
	c.Send(response)
}

func GetCharacterById(c *fiber.Ctx) {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	character, err := findOne(c, filter)
	if err {
		return
	}

	response, _ := json.Marshal(character)
	c.Send(response)
}

func GetCharacterByName(c *fiber.Ctx) {
	name := c.Params("name")
	filter := bson.M{"name": name}

	character, err := findOne(c, filter)
	if err {
		return
	}

	response, _ := json.Marshal(character)
	c.Send(response)
}

func AddCharacter(c *fiber.Ctx) {
	var character document.Character

	// we decode our body request params
	json.Unmarshal([]byte(c.Body()), &character)

	// insert our character model.
	result, err := collection.InsertOne(context.Background(), character)
	if err != nil {
		database.GetError(err, c)
		return
	}

	response, _ := json.Marshal(result)
	c.Status(fiber.StatusCreated).Send(response)
}

func UpdateCharacter(c *fiber.Ctx) {
	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	// Create filter
	filter := bson.M{"_id": id}

	var character document.Character

	// Read update model from body request
	json.Unmarshal([]byte(c.Body()), &character)

	update := bson.M{
		"$set": character,
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&character)

	if err != nil {
		database.GetError(err, c)
		return
	}

	character.Id = id.Hex()
	response, _ := json.Marshal(character)
	c.Send(response)
}

func PartialUpdateCharacter(c *fiber.Ctx) {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{"_id": id}

	dbCharacter, err := findOne(c, filter)
	if err {
		return
	}

	var character document.Character

	// Read update model from body request
	json.Unmarshal([]byte(c.Body()), &character)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", nullIf(character.Name, dbCharacter.Name)},
			{"description", nullIf(character.Description, dbCharacter.Description)},
			{"superPowers", nullIf(character.SuperPowers, dbCharacter.SuperPowers)},
		}},
	}

	res, err2 := collection.UpdateOne(context.Background(), filter, update)

	if err2 != nil {
		database.GetError(err2, c)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func DeleteCharacter(c *fiber.Ctx) {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		database.GetError(err, c)
		return
	}

	jsonResponse, _ := json.Marshal(res)
	c.Send(jsonResponse)
}

func findOne(c *fiber.Ctx, filter bson.M) (document.Character, bool) {
	var character document.Character

	err := collection.FindOne(context.Background(), filter).Decode(&character)

	if err != nil {
		database.GetError(err, c)
		return document.Character{}, true
	}

	return character, false
}

func nullIf(s1 string, s2 string) string {
	if s1 != "" {
		return s1
	} else {
		return s2
	}
}
