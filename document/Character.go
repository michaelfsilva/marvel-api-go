package document

import "go.mongodb.org/mongo-driver/bson/primitive"

// the fields between backpicks define the db mapping, for mongo is bson
type Character struct {
	// Id          string `json:"_id" bson:"_id,omitempty"`
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	SuperPowers string             `json:"superPowers" bson:"superPowers,omitempty"`
}
