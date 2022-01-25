package repository

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"marvel-api-go/database"
	"marvel-api-go/document"
)

//Connection mongoDB
var collection = database.GetCollection("character")

func ListAll(c *fiber.Ctx) []document.Character {
	var characters []document.Character

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.Background(), bson.M{})

	// Close the cursor once finished
	// A defer statement defers the execution of a function until the surrounding function returns.
	// simply, run cur.Close() process but after cur.Next() finished.
	defer func(cur *mongo.Cursor, ctx context.Context) {
		cur.Close(ctx)
	}(cur, context.Background())

	if err != nil {
		database.GetError(err, c)
		return nil
	}

	// better than using a loop
	err = cur.All(context.Background(), &characters)
	if err != nil {
		return nil
	}

	return characters
}

func GetById(c *fiber.Ctx) document.Character {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	character, err := findOne(c, filter)
	if err {
		return document.Character{}
	}

	return character
}

func GetByName(c *fiber.Ctx) []document.Character {
	name := c.Params("name")
	filter := bson.M{"name": name}

	var characters []document.Character

	cur, err := collection.Find(context.Background(), filter)

	defer func(cur *mongo.Cursor, ctx context.Context) {
		cur.Close(ctx)
	}(cur, context.Background())

	if err != nil {
		database.GetError(err, c)
		return nil
	}

	err = cur.All(context.Background(), &characters)
	if err != nil {
		return nil
	}

	return characters
}

func Add(c *fiber.Ctx) *mongo.InsertOneResult {
	var character document.Character

	// we decode our body request params
	json.Unmarshal(c.Body(), &character)

	// insert our character model.
	result, err := collection.InsertOne(context.Background(), character)
	if err != nil {
		database.GetError(err, c)
		return nil
	}

	return result
}

func Update(c *fiber.Ctx) document.Character {
	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	// Create filter
	filter := bson.M{"_id": id}

	var character document.Character

	// Read update model from body request
	json.Unmarshal(c.Body(), &character)

	update := bson.M{
		"$set": character,
	}

	result := collection.FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() != nil {
		return document.Character{}
	}

	character.Id = id.Hex()
	return character
}

func PartialUpdate(c *fiber.Ctx) *mongo.UpdateResult {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{"_id": id}

	dbCharacter, err := findOne(c, filter)
	if err {
		return nil
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

	response, err2 := collection.UpdateOne(context.Background(), filter, update)
	if err2 != nil {
		database.GetError(err2, c)
		return nil
	}

	return response
}

func Delete(c *fiber.Ctx) *mongo.DeleteResult {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	response, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		database.GetError(err, c)
		return nil
	}

	return response
}

func findOne(c *fiber.Ctx, filter bson.M) (document.Character, bool) {
	var character document.Character

	err := collection.FindOne(context.Background(), filter).Decode(&character)

	if err != nil {
		database.GetErrorWithStatus(err, c, fiber.StatusNotFound)
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
