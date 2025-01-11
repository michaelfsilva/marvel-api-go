package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"marvel-api-go/database"
	"marvel-api-go/document"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharacterRepository interface {
	ListAll() []document.Character
}

type CharacterRepositoryImpl struct {
	database database.Database
}

func NewCharacterRepository(database database.Database) *CharacterRepositoryImpl {
	return &CharacterRepositoryImpl{database}
}

func (r *CharacterRepositoryImpl) ListAll(c *fiber.Ctx) []document.Character {
	var characters []document.Character

	// TODO ta retornando a collection inteira ou é uma conexão?
	println(r.database.Collection)
	fmt.Println(r.database.Collection)

	// bson.M{},  we passed empty filter. So we want to get all data.
	cursor, err := r.database.Collection.Find(context.Background(), bson.M{})

	// Close the cursor once finished
	// A defer statement defers the execution of a function until the surrounding function returns.
	// simply, run cur.Close() process but after cur.Next() finished.
	defer func(cur *mongo.Cursor, ctx context.Context) {
		cur.Close(ctx)
	}(cursor, context.Background())

	if err != nil {
		database.GetError(err, c)
		return nil
	}

	// better than using a loop
	err = cursor.All(context.Background(), &characters)
	if err != nil {
		return nil
	}

	return characters
}

func (r *CharacterRepositoryImpl) GetById(c *fiber.Ctx) document.Character {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	character, err := r.findOne(c, filter)
	if err {
		return document.Character{}
	}

	return character
}

func (r *CharacterRepositoryImpl) GetByName(c *fiber.Ctx) []document.Character {
	name := c.Params("name")
	filter := bson.M{"name": name}

	var characters []document.Character

	cur, err := r.database.Collection.Find(context.Background(), filter)

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

func (r *CharacterRepositoryImpl) Add(c *fiber.Ctx) document.Character {
	var character document.Character

	// we decode our body request params
	json.Unmarshal(c.Body(), &character)

	// insert our character model.
	result, err := r.database.Collection.InsertOne(context.Background(), character)
	if err != nil {
		database.GetError(err, c)
		return document.Character{}
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		character.Id = id.Hex()
	} /*  else {
		character.Id = id
	} */

	return character
}

func (r *CharacterRepositoryImpl) Update(c *fiber.Ctx) document.Character {
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

	result := r.database.Collection.FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() != nil {
		return document.Character{}
	}

	character.Id = id.Hex()
	return character
}

func (r *CharacterRepositoryImpl) PartialUpdate(c *fiber.Ctx) document.Character {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{"_id": id}

	dbCharacter, err := r.findOne(c, filter)
	if err {
		return document.Character{}
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

	_, err2 := r.database.Collection.UpdateOne(context.Background(), filter, update)
	if err2 != nil {
		database.GetError(err2, c)
		return document.Character{}
	}

	character.Id = id.Hex()
	return character
}

func (r *CharacterRepositoryImpl) Delete(c *fiber.Ctx) string {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	_, err := r.database.Collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		database.GetError(err, c)
		return ""
	}

	return id.Hex()
}

func (r *CharacterRepositoryImpl) findOne(c *fiber.Ctx, filter bson.M) (document.Character, bool) {
	var character document.Character

	err := r.database.Collection.FindOne(context.Background(), filter).Decode(&character)

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
