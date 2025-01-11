package service

import (
	"testing"

	. "marvel-api-go/document"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// func (m *MockRepository) FetchData() (string, error) {
// 	args := m.Called()
// 	return args.String(0), args.Error(1)
// }

func (m *MockRepository) ListAll() []Character {
	args := m.Called()
	return args.Get(0).([]Character)
}

var mockRepo = new(MockRepository)
var service = NewCharacterService(mockRepo)

func TestShouldListAllCharacters(t *testing.T) {
	// mockRepo.On("FetchData").Return("Dados Mockados", nil)
	mockRepo.On("ListAll").Return([]Character{
		{Id: "1", Name: "Character 1"},
		{Id: "2", Name: "Character 2"},
	})

	result := service.ListAll(new(fiber.Ctx))

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "1", result[0].Id)
	assert.Equal(t, "Character 1", result[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestShouldGetCharacterById(t *testing.T) {
}

func TestShouldGetCharacterByName(t *testing.T) {
}

func TestShouldAddCharacter(t *testing.T) {
}

func TestShouldUpdateCharacter(t *testing.T) {
}

func TestShouldPartialUpdateCharacter(t *testing.T) {
}

func TestShouldDeleteCharacter(t *testing.T) {}
