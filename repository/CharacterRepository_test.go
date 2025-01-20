package repository

import (
	"marvel-api-go/database"
	"marvel-api-go/document"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestListAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		database.Collection = mt.Coll
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "john"},
			{"description", "test1"},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"name", "john"},
			{"description", "test2"},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		characters, err := ListAll()

		assert.Nil(t, err)
		assert.Equal(t, []document.Character{
			{ID: id1, Name: "john", Description: "test1"},
			{ID: id2, Name: "john", Description: "test2"},
		}, characters)
	})
}

func TestGetById(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		database.Collection = mt.Coll
		expectedCharacter := document.Character{
			ID:          primitive.NewObjectID(),
			Name:        "john",
			Description: "test",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedCharacter.ID},
			{"name", expectedCharacter.Name},
			{"description", expectedCharacter.Description},
		}))

		result, err := GetById(expectedCharacter.ID.Hex())

		assert.Nil(t, err)
		assert.Equal(t, &expectedCharacter, result)
	})
}

func TestGetByName(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		database.Collection = mt.Coll
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "john"},
			{"description", "test1"},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"name", "john"},
			{"description", "test2"},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		characters, err := GetByName("john")

		assert.Nil(t, err)
		assert.Equal(t, []document.Character{
			{ID: id1, Name: "john", Description: "test1"},
			{ID: id2, Name: "john", Description: "test2"},
		}, characters)
	})
}

func TestAdd(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		database.Collection = mt.Coll
		id := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		insertedCharacter, err := Add(document.Character{
			ID:          id,
			Name:        "john",
			Description: "test",
		})

		assert.Nil(t, err)
		assert.Equal(t, &document.Character{
			ID:          id,
			Name:        "john",
			Description: "test",
		}, insertedCharacter)
	})

	mt.Run("custom error duplicate", func(mt *mtest.T) {
		database.Collection = mt.Coll
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		insertedCharacter, err := Add(document.Character{})

		assert.Nil(t, insertedCharacter)
		assert.NotNil(t, err)
		assert.True(t, mongo.IsDuplicateKeyError(err))
	})

	mt.Run("simple error", func(mt *mtest.T) {
		database.Collection = mt.Coll
		mt.AddMockResponses(bson.D{{"ok", 0}})

		insertedCharacter, err := Add(document.Character{})

		assert.Nil(t, insertedCharacter)
		assert.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		database.Collection = mt.Coll
		characterData := document.Character{
			ID:          primitive.NewObjectID(),
			Name:        "john",
			Description: "test",
		}
		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"_id", characterData.ID},
				{"name", characterData.Name},
				{"description", characterData.Description},
			}},
		})

		updatedCharacter, err := Update(characterData)

		assert.Nil(t, err)
		assert.Equal(t, &characterData, updatedCharacter)
	})
}

// TODO fix this test/code
// func TestPartialUpdate(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("success", func(mt *mtest.T) {
// 		database.Collection = mt.Coll
// 		characterData := document.Character{
// 			ID:          primitive.NewObjectID(),
// 			Description: "test2",
// 		}
// 		mt.AddMockResponses(bson.D{
// 			{"ok", 1},
// 			{"value", bson.D{
// 				{"_id", characterData.ID},
// 				{"name", "john"},
// 				{"description", "test"},
// 			}},
// 		})

// 		updatedCharacter, err := PartialUpdate(characterData)

// 		assert.Nil(t, err)
// 		assert.Equal(t, "john", updatedCharacter.Name)
// 		assert.Equal(t, &characterData.Description, updatedCharacter.Description)
// 	})
// }

func TestDeleteOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		database.Collection = mt.Coll
		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})
		err := Delete(primitive.NewObjectID())
		assert.Nil(t, err)
	})

	mt.Run("no document deleted", func(mt *mtest.T) {
		database.Collection = mt.Coll
		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 0}})
		err := Delete(primitive.NewObjectID())
		assert.NotNil(t, err)
	})
}
