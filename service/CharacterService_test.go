package service

import (
	"testing"

	. "marvel-api-go/document"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mockRepo = new(MockRepository)
var service = NewCharacterService(mockRepo)

func TestShouldListAllCharacters(t *testing.T) {
	id1 := primitive.NewObjectID()
	mockRepo.On("ListAll").Return(
		[]Character{
			{ID: id1, Name: "Character 1"},
			{ID: primitive.NewObjectID(), Name: "Character 2"},
		},
		nil,
	)

	result, _ := service.ListAll()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, id1, result[0].ID)
	assert.Equal(t, "Character 1", result[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestShouldGetCharacterById(t *testing.T) {
	id1 := primitive.NewObjectID()
	characterMock := Character{
		ID: id1, Name: "Character 1", Description: "", SuperPowers: "",
	}
	mockRepo.On("GetById", mock.Anything).Return(&characterMock, nil)

	result, _ := service.GetCharacterById("1")

	assert.Equal(t, id1, result.ID)
	assert.Equal(t, "Character 1", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestShouldGetCharacterByName(t *testing.T) {
	id1 := primitive.NewObjectID()
	mockRepo.On("GetByName", mock.Anything).Return(
		[]Character{
			{ID: id1, Name: "Character 1"},
			{ID: primitive.NewObjectID(), Name: "Character 2"},
		},
		nil,
	)

	result, _ := service.GetCharacterByName("Character")

	assert.Equal(t, 2, len(result))
	assert.Equal(t, id1, result[0].ID)
	assert.Equal(t, "Character 1", result[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestShouldAddCharacter(t *testing.T) {
	id1 := primitive.NewObjectID()
	characterMock := Character{
		ID: id1, Name: "Character 1", Description: "", SuperPowers: "",
	}
	mockRepo.On("Add", mock.Anything).Return(&characterMock, nil)

	result, _ := service.AddCharacter(
		Character{
			Name: "Character 1", Description: "", SuperPowers: "",
		},
	)

	assert.Equal(t, id1, result.ID)
	assert.Equal(t, "Character 1", result.Name)

	mockRepo.AssertExpectations(t)
}

func TestShouldUpdateCharacter(t *testing.T) {
	id1 := primitive.NewObjectID()
	characterMock := Character{
		ID: id1, Name: "Character 1", Description: "", SuperPowers: "",
	}
	mockRepo.On("Update", mock.Anything).Return(&characterMock, nil)

	result, _ := service.UpdateCharacter(
		Character{
			Name: "Character 1", Description: "", SuperPowers: "",
		},
	)

	assert.Equal(t, id1, result.ID)
	assert.Equal(t, "Character 1", result.Name)

	mockRepo.AssertExpectations(t)
}

// func TestShouldPartialUpdateCharacter(t *testing.T) {
// 	id1 := primitive.NewObjectID()
// 	characterMock := Character{
// 		ID: id1, Name: "Character 1", Description: "", SuperPowers: "",
// 	}
// 	mockRepo.On("Update", mock.Anything).Return(&characterMock, nil)

// 	result, _ := service.UpdateCharacter(
// 		Character{
// 			Name: "Character 1", Description: "", SuperPowers: "",
// 		},
// 	)

// 	assert.Equal(t, id1, result.ID)
// 	assert.Equal(t, "Character 1", result.Name)

// 	mockRepo.AssertExpectations(t)
// }

func TestShouldDeleteCharacter(t *testing.T) {
	mockRepo.On("Delete", mock.Anything).Return(nil)

	service.DeleteCharacter("1")

	mockRepo.AssertExpectations(t)
}
