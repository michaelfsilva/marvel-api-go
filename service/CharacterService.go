package service

import (
	"marvel-api-go/document"
	r "marvel-api-go/repository"

	"github.com/gofiber/fiber/v2"
)

type CharacterService struct {
	CharacterRepository r.CharacterRepository
}

func NewCharacterService(repository r.CharacterRepository) *CharacterService {
	return &CharacterService{CharacterRepository: repository}
}

func (s *CharacterService) ListAll(c *fiber.Ctx) []document.Character {
	return s.CharacterRepository.ListAll()
}

// var ListAll = func(c *fiber.Ctx) []document.Character {
// 	return repository.ListAll(c)
// }

// func GetCharacterById(c *fiber.Ctx) document.Character {
// 	return repository.GetById(c)
// }

// func GetCharacterByName(c *fiber.Ctx) []document.Character {
// 	return repository.GetByName(c)
// }

// func AddCharacter(c *fiber.Ctx) document.Character {
// 	return repository.Add(c)
// }

// func UpdateCharacter(c *fiber.Ctx) document.Character {
// 	return repository.Update(c)
// }

// func PartialUpdateCharacter(c *fiber.Ctx) document.Character {
// 	return repository.PartialUpdate(c)
// }

// func DeleteCharacter(c *fiber.Ctx) string {
// 	return repository.Delete(c)
// }
