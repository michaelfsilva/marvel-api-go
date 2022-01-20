package document

type Character struct {
	Id          string `json:"id" bson:"id,omitempty"`
	Name        string `json:"name" bson:"name,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
	SuperPowers string `json:"superPowers" bson:"superPowers,omitempty"`
}
