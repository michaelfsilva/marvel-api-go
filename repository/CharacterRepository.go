package repository

import (
	"context"
	"errors"
	"marvel-api-go/database"
	"marvel-api-go/document"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharacterRepository interface {
	ListAll() ([]document.Character, error)
	GetById(id string) (*document.Character, error)
	GetByName(name string) ([]document.Character, error)
	Add(character document.Character) (*document.Character, error)
	Update(character document.Character) (*document.Character, error)
	PartialUpdate(character document.Character) (document.Character, error)
	Delete(id string) error
}

type CharacterRepositoryImpl struct{}

func (r *CharacterRepositoryImpl) InitRepository(connectionString string) {
	database.NewDatabase(connectionString, "character")
}

func (r *CharacterRepositoryImpl) ListAll() ([]document.Character, error) {
	var characters []document.Character

	// bson.M{},  we passed empty filter. So we want to get all data.
	cursor, err := database.Collection.Find(context.Background(), bson.M{})

	if err != nil {
		// database.GetError(err, c) // TODO
		return nil, err
	}

	// Close the cursor once finished
	// A defer statement defers the execution of a function until the surrounding function returns.
	// simply, run cur.Close() process but after cur.Next() finished.
	defer func(cur *mongo.Cursor, ctx context.Context) {
		cur.Close(ctx)
	}(cursor, context.Background())

	// better than using a loop
	err = cursor.All(context.Background(), &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharacterRepositoryImpl) GetById(id string) (*document.Character, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	var character document.Character

	if err := database.Collection.FindOne(context.Background(), filter).Decode(&character); err != nil {
		if err == mongo.ErrNoDocuments {
			return &character, nil
		}
		return nil, err
	}

	return &character, nil
}

func (r *CharacterRepositoryImpl) GetByName(name string) ([]document.Character, error) {
	filter := bson.M{"name": name}
	var characters []document.Character

	cur, err := database.Collection.Find(context.Background(), filter)
	if err != nil {
		// database.GetError(err, c)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		cur.Close(ctx)
	}(cur, context.Background())

	err = cur.All(context.Background(), &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharacterRepositoryImpl) Add(character document.Character) (*document.Character, error) {
	result, err := database.Collection.InsertOne(context.Background(), character)
	if err != nil {
		// database.GetError(err, c) // TODO
		return nil, err
	}

	character.ID = result.InsertedID.(primitive.ObjectID)
	return &character, nil
}

func (r *CharacterRepositoryImpl) Update(character document.Character) (*document.Character, error) {
	// Create filter
	filter := bson.M{"_id": character.ID}
	update := bson.M{"$set": character}

	err := database.Collection.FindOneAndUpdate(
		context.Background(), filter, update, //options.FindOneAndUpdate().SetReturnDocument(1),
	).Decode(&character)
	if err != nil {
		return nil, err
	}

	return &character, nil
}

func (r *CharacterRepositoryImpl) PartialUpdate(character document.Character) (document.Character, error) {
	filter := bson.M{"_id": character.ID}

	dbCharacter, err := r.GetById(character.ID.Hex())
	if err != nil {
		return document.Character{}, err
	}

	// prepare update model
	update := bson.D{
		{"$set", bson.D{
			{"name", nullIf(character.Name, dbCharacter.Name)},
			{"description", nullIf(character.Description, dbCharacter.Description)},
			{"superPowers", nullIf(character.SuperPowers, dbCharacter.SuperPowers)},
		}},
	}

	_, err = database.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// database.GetError(err2, c) // TODO
		return document.Character{}, err
	}

	return character, nil
}

func (r *CharacterRepositoryImpl) Delete(id string) error {
	result, err := database.Collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if err != nil {
		// database.GetError(err, c) // TODO
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document deleted")
	}

	return nil
}

func nullIf(s1 string, s2 string) string {
	if s1 != "" {
		return s1
	} else {
		return s2
	}
}
