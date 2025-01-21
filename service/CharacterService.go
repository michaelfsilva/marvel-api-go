package service

import (
	"marvel-api-go/document"
	r "marvel-api-go/repository"
)

type CharacterService struct {
	characterRepository r.CharacterRepository
}

func NewCharacterService(repository r.CharacterRepository) *CharacterService {
	return &CharacterService{characterRepository: repository}
}

func (s *CharacterService) ListAll() ([]document.Character, error) {
	return s.characterRepository.ListAll()
}

func (s *CharacterService) GetCharacterById(id string) (*document.Character, error) {
	return s.characterRepository.GetById(id)
}

func (s *CharacterService) GetCharacterByName(name string) ([]document.Character, error) {
	return s.characterRepository.GetByName(name)
}

func (s *CharacterService) AddCharacter(character document.Character) (*document.Character, error) {
	return s.characterRepository.Add(character)
}

func (s *CharacterService) UpdateCharacter(character document.Character) (*document.Character, error) {
	return s.characterRepository.Update(character)
}

func (s *CharacterService) PartialUpdateCharacter(character document.Character) (document.Character, error) {
	return s.characterRepository.PartialUpdate(character)
}

func (s *CharacterService) DeleteCharacter(id string) error {
	return s.characterRepository.Delete(id)
}
