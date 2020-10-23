package domain

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User domain model
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Username  string             `json:"username" validate:"required"`
	Email     string             `json:"email" validate:"email"`
	Salt      string             `json:"salt"`
	Hash      string             `json:"hash"`
	CreatedOn string             `json:"-"`
	UpdatedOn string             `json:"-"`
	DeletedOn string             `json:"-"`
}

// FromJSON decodes JSON data
func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// ToJSON serializes the contents of the collection to JSON
func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// Validate User domain model
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
