package repository

import (
	"context"
	"marvel-api-go/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func TestDatabase(t *testing.T) {
	// Create a new instance of MockMongoCollection
	mockCollection := new(MockCollection)

	// Define mock behavior for InsertOne
	mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{
		InsertedID: "1234",
	}, nil)

	// Define mock behavior for FindOne
	mockCollection.On("FindOne", mock.Anything, mock.Anything).Return(&mongo.SingleResult{})

	// Create the original Database instance with the mock collection
	db := &database.Database{
		Collection: mockCollection, // Inject the mock collection
	}

	// Test AddItem method
	item := bson.M{"name": "test item"}
	result, err := db.AddItem(item)
	assert.NoError(t, err)
	assert.Equal(t, "1234", result.InsertedID)

	// Verify that InsertOne was called as expected
	mockCollection.AssertExpectations(t)

	// Test GetItem method
	filter := bson.M{"name": "test item"}
	result2 := db.GetItem(filter)

	// Verify that FindOne was called as expected
	mockCollection.AssertExpectations(t)
}
