# marvel-api-go
Marvel API using Golang

## Instructions to use
Download the zip project by clicking on Clone or download button, then Download ZIP;  
Extract the zip file and open a terminal in the folder extracted;  
Run the following commands:

```
$ go run marvel-api-go.go
```

Access http://localhost:8080/api/characters/ with user "test" and password "marvel" to see a list of characters.

## APIs endpoints
GET http://localhost:8080/api/characters/ [list all characters]  
GET http://localhost:8080/api/characters/{id} [list a character by ID]  
GET http://localhost:8080/api/characters/findByName/{name} [list a character by name ignoring case]  
POST http://localhost:8080/api/characters/ [add a new character]  
PUT http://localhost:8080/api/characters/{id} [update character all attributes]  
PATCH http://localhost:8080/api/characters/{id} [update one or more attributes of a character]  
DELETE http://localhost:8080/api/characters/{id} [remove a character]
