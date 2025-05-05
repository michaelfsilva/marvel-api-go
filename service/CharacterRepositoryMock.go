package service

import (
	"marvel-api-go/document"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListAll() ([]document.Character, error) {
	args := m.Called()
	return args.Get(0).([]document.Character), args.Error(1)
}

func (m *MockRepository) GetById(id string) (*document.Character, error) {
	args := m.Called(id)
	return args.Get(0).(*document.Character), args.Error(1)
}

func (m *MockRepository) GetByName(name string) ([]document.Character, error) {
	args := m.Called(name)
	return args.Get(0).([]document.Character), args.Error(1)
}

func (m *MockRepository) Add(character document.Character) (*document.Character, error) {
	args := m.Called(character)
	return args.Get(0).(*document.Character), args.Error(1)
}

func (m *MockRepository) Update(character document.Character) (document.Character, error) {
	args := m.Called(character)
	return args.Get(0).(document.Character), args.Error(1)
}

func (m *MockRepository) PartialUpdate(character document.Character) (document.Character, error) {
	args := m.Called(character)
	return args.Get(0).(document.Character), args.Error(1)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
