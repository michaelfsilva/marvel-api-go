package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var testEmptyHandler = func(c *fiber.Ctx) error {
	return nil
}

func TestShouldReturn200WhenGetAllIsCalled(t *testing.T) {
	dummyHandler := testEmptyHandler

	app := fiber.New()

	app.Get("/api/characters123", dummyHandler)
	testStatus(t, app, "/api/characters123", fiber.MethodGet, fiber.StatusOK)
}

func TestShouldReturn204WhenGetByIdWithWrongID(t *testing.T) {
	dummyHandler := testEmptyHandler

	app := fiber.New()

	app.Get("/api/characters/1234", dummyHandler)
	testStatus(t, app, "/api/characters", fiber.MethodGet, fiber.StatusNoContent)
}

func TestShouldReturn201WhenAddIsCalled(t *testing.T) {
	dummyHandler := testEmptyHandler

	app := fiber.New()

	app.Post("/api/characters", dummyHandler)
	testStatus(t, app, "/api/characters", fiber.MethodPost, fiber.StatusCreated)
}

func TestShouldReturn404WhenPutIsCalledWithWrongID(t *testing.T) {
	dummyHandler := testEmptyHandler

	app := fiber.New()

	app.Put("/api/characters", dummyHandler)
	testStatus(t, app, "/api/characters", fiber.MethodPut, fiber.StatusNotFound)
}

func TestShouldReturn404WhenPatchIsCalledWithWrongID(t *testing.T) {
	dummyHandler := testEmptyHandler

	app := fiber.New()

	app.Patch("/api/characters", dummyHandler)
	testStatus(t, app, "/api/characters", fiber.MethodPatch, fiber.StatusNotFound)
}

func TestShouldReturn404WhenDeleteIsCalledWithWrongID(t *testing.T) {
	dummyHandler := testEmptyHandler

	app := fiber.New()

	app.Delete("/api/characters", dummyHandler)
	testStatus(t, app, "/api/characters", fiber.MethodDelete, fiber.StatusNotFound)
}

func testStatus(t *testing.T, app *fiber.App, url string, method string, statusCode int) {
	t.Helper()

	req := httptest.NewRequest(method, url, nil)

	resp, err := app.Test(req)
	assert.Equal(t, nil, err, "app.Test(req)")
	assert.Equal(t, statusCode, resp.StatusCode, "Status code")
}
