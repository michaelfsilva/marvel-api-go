package controller

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"marvel-api-go/database"
	"marvel-api-go/document"
)

//Connection mongoDB
var collection = database.GetCollection("character")

func GetCharacters(c *fiber.Ctx) { //(w http.ResponseWriter, r *http.Request) {
	var characters []document.Character

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		database.GetError(err, w)
		return
	}

	// Close the cursor once finished
	// A defer statement defers the execution of a function until the surrounding function returns.
	// simply, run cur.Close() process but after cur.Next() finished.
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var character document.Character
		// & character returns the memory address of the following variable.
		err := cur.Decode(&character) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		characters = append(characters, character)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(characters) // encode similar to serialize process.
}

func GetCharacterById(c *fiber.Ctx) { //(w http.ResponseWriter, r *http.Request) {
	var filter bson.M = bson.M{}

	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M

	//var character document.Character

	// we get params with mux.
	//var params = mux.Vars(r)

	// string to primitive.ObjectID
	//id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	//filter := bson.M{"_id": id}
	//err := collection.FindOne(context.TODO(), filter).Decode(&character)
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		//database.GetError(err, w)
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	//json.NewEncoder(w).Encode(character)
	json, _ := json.Marshal(results)
	c.Send(json)
}

func GetCharacterByName(c *fiber.Ctx) { //(w http.ResponseWriter, r *http.Request) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var filter bson.M = bson.M{}

	if c.Params("name") != "" {
		name := c.Params("name")
		objName, _ := primitive.ObjectIDFromHex(name)
		filter = bson.M{"_name": objName}
	}

	var results []bson.M

	//var character document.Character

	// we get params with mux.
	//var params = mux.Vars(r)

	// string to primitive.ObjectID
	//name, _ := primitive.ObjectIDFromHex(params["name"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	//filter := bson.M{"_name": name}
	//err := collection.FindOne(context.TODO(), filter).Decode(&character)
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		//database.GetError(err, w)
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	//json.NewEncoder(w).Encode(character)
	json, _ := json.Marshal(results)
	c.Send(json)
}

func AddCharacter(c *fiber.Ctx) { //(w http.ResponseWriter, r *http.Request) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var character document.Character
	json.Unmarshal([]byte(c.Body()), &character)

	// we decode our body request params
	//_ = json.NewDecoder(r.Body).Decode(&character)

	// insert our character model.
	//result, err := collection.InsertOne(context.TODO(), character)
	result, err := collection.InsertOne(context.Background(), character)
	if err != nil {
		//database.GetError(err, w)
		c.Status(500).Send(err)
		return
	}

	//json.NewEncoder(w).Encode(result)
	response, _ := json.Marshal(result)
	c.Send(response)
}

func UpdateCharacter(c *fiber.Ctx) { //(w http.ResponseWriter, r *http.Request) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	//var params = mux.Vars(r)

	//Get id from parameters
	//id, _ := primitive.ObjectIDFromHex(params["id"])

	var character document.Character
	json.Unmarshal([]byte(c.Body()), &character)

	// Create filter
	//filter := bson.M{"_id": id}

	// Read update model from body request
	//_ = json.NewDecoder(r.Body).Decode(&character)

	// prepare update model.
	//update := bson.D{
	//	{"$set", bson.D{
	//		{"name", character.Name},
	//		{"description", character.Description},
	//		{"superPowers", character.SuperPowers},
	//	}},
	//}

	update := bson.M{
		"$set": character,
	}

	//err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&character)
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		//database.GetError(err, w)
		c.Status(500).Send(err)
		return
	}

	//character.Id = id.Hex()

	//json.NewEncoder(w).Encode(character)

	response, _ := json.Marshal(res)
	c.Send(response)
}

func DeleteCharacter(c *fiber.Ctx) { //(w http.ResponseWriter, r *http.Request) {
	collection, err := database.GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	//var params = mux.Vars(r)

	// string to primitve.ObjectID
	//id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	//filter := bson.M{"_id": id}

	//deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		//database.GetError(err, w)
		c.Status(500).Send(err)
		return
	}

	//json.NewEncoder(w).Encode(deleteResult)
	jsonResponse, _ := json.Marshal(res)
	c.Send(jsonResponse)
}
