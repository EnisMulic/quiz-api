package domain

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

// User domain model
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required, unique"`
	Email     string `json:"email" validate:"email, required, unique"`
	Salt      string `json:"salt"`
	Hash      string `json:"hash"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// FromJSON decodes JSON data
func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// Validate User domain model
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
