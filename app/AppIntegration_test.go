package app

import (
	"bytes"
	"context"
	"fmt"
	. "marvel-api-go/document"
	"marvel-api-go/repository"

	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOnlyRepository(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(
		getMongoContainerConnection(t),
	)

	id1 := primitive.NewObjectID()
	doc := Character{id1, "Test", "", ""}
	repository.Add(doc)

	result, err := repository.GetByName("Test")

	assert.Nil(t, err)
	assert.Equal(t, []Character{{id1, "Test", "", ""}}, result)
}

func TestShouldReturn200WhenGetAllIsCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	doc := Character{id1, "Test", "", ""}
	doc2 := Character{id2, "Test2", "", ""}
	repository.Add(doc)
	repository.Add(doc2)

	req := httptest.NewRequest("GET", "/api/characters", nil)
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var characters []Character
	err = json.NewDecoder(resp.Body).Decode(&characters)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	assert.Equal(t, 2, len(characters))
	assert.Equal(t, []Character{doc, doc2}, characters)
}

func TestShouldReturn200WhenGetByIdIsCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	doc := Character{id1, "Test", "", ""}
	doc2 := Character{id2, "Test2", "", ""}
	repository.Add(doc)
	repository.Add(doc2)

	req := httptest.NewRequest("GET", "/api/characters/"+id1.Hex(), nil)
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var character Character
	err = json.NewDecoder(resp.Body).Decode(&character)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	assert.Equal(t, doc, character)
}

func TestShouldReturn204WhenGetByIdWithWrongID(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id1 := primitive.NewObjectID()
	doc := Character{id1, "Test", "", ""}
	repository.Add(doc)

	req := httptest.NewRequest("GET", "/api/characters/1234", nil)
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestShouldReturn200WhenGetByNameCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	doc := Character{primitive.NewObjectID(), "Test", "", ""}
	doc2 := Character{primitive.NewObjectID(), "Test2", "", ""}
	repository.Add(doc)
	repository.Add(doc2)

	req := httptest.NewRequest("GET", "/api/characters/findByName/"+"Test", nil)
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var characters []Character
	err = json.NewDecoder(resp.Body).Decode(&characters)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	assert.Equal(t, []Character{doc}, characters)
}

func TestShouldReturn204WhenGetByNameWithWrongName(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	doc := Character{primitive.NewObjectID(), "Test", "", ""}
	repository.Add(doc)

	req := httptest.NewRequest("GET", "/api/characters/findByName/"+"abc", nil)
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestShouldReturn201WhenAddIsCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	character := Character{primitive.NewObjectID(), "Test", "", ""}

	body, _ := json.Marshal(character)

	req := httptest.NewRequest("POST", "/api/characters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result Character
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error while decoding response: %v", err)
	}

	assert.Equal(t, character, result)
}

func TestShouldReturn200WhenPutIsCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id := primitive.NewObjectID()
	repository.Add(Character{id, "Test", "", ""})

	character := Character{id, "Name", "Test", "Test"}
	body, _ := json.Marshal(character)

	req := httptest.NewRequest("PUT", "/api/characters/"+id.Hex(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result Character
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error while decoding response: %v", err)
	}

	assert.Equal(t, character, result)

	dbResult, _ := repository.GetById(id.Hex())
	assert.Equal(t, character, *dbResult)
}

func TestShouldReturn404WhenPutIsCalledWithWrongID(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id := primitive.NewObjectID()
	repository.Add(Character{id, "Test", "", ""})

	character := Character{id, "Name", "Test", "Test"}
	body, _ := json.Marshal(character)

	req := httptest.NewRequest("PUT", "/api/characters/1234", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestShouldReturn200WhenPatchIsCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id := primitive.NewObjectID()
	repository.Add(Character{id, "Test", "Test", ""})

	character := Character{id, "Name", "", "Test"}
	body, _ := json.Marshal(character)

	req := httptest.NewRequest("PATCH", "/api/characters/"+id.Hex(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result Character
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Error while decoding response: %v", err)
	}

	expectedResult := Character{id, "Name", "Test", "Test"}
	assert.Equal(t, expectedResult, result)

	dbResult, _ := repository.GetById(id.Hex())
	assert.Equal(t, expectedResult, *dbResult)
}

func TestShouldReturn404WhenPatchIsCalledWithWrongID(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id := primitive.NewObjectID()
	character := Character{id, "Test", "Test", ""}
	repository.Add(character)

	body, _ := json.Marshal(Character{id, "Name", "", "Test"})

	req := httptest.NewRequest("PATCH", "/api/characters/"+"1234", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	dbResult, _ := repository.GetById(id.Hex())
	assert.Equal(t, character, *dbResult)
}

func TestShouldReturn200WhenDeleteIsCalled(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id := primitive.NewObjectID()
	character := Character{id, "Test", "Test", ""}
	repository.Add(character)

	req := httptest.NewRequest("DELETE", "/api/characters/"+id.Hex(), nil)
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	dbResult, _ := repository.GetById(id.Hex())
	assert.Nil(t, dbResult)
}

func TestShouldReturn404WhenDeleteIsCalledWithWrongID(t *testing.T) {
	repository := repository.CharacterRepositoryImpl{}
	repository.InitRepository(getMongoContainerConnection(t))
	app := SetupApp(&repository)

	id := primitive.NewObjectID()
	character := Character{id, "Test", "Test", ""}
	repository.Add(character)

	req := httptest.NewRequest("DELETE", "/api/characters/"+"1234", nil)
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	dbResult, _ := repository.GetById(id.Hex())
	assert.Equal(t, character, *dbResult)
}

func setAuthHeader(req *http.Request) {
	username := "test"
	password := "marvel"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	req.Header.Set("Authorization", auth)
}

func getMongoContainerConnection(t *testing.T) string {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	mongoContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Error starting MongoDB container: %v", err)
	}

	// TODO check this
	// defer mongoContainer.Terminate(context.Background())

	host, err := mongoContainer.Host(context.Background())
	if err != nil {
		t.Fatalf("Error getting container host: %v", err)
	}
	port, err := mongoContainer.MappedPort(context.Background(), "27017")
	if err != nil {
		t.Fatalf("Error getting mapped port: %v", err)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port.Port())
}
