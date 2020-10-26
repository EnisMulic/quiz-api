package domain

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User domain model
type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `json:"username" validate:"required"`
	Email    string             `json:"email" validate:"email"`
	Salt     string             `json:"salt"`
	Hash     string             `json:"hash"`
}

// Validate User domain model
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// UserUpsertRequest struct
type UserUpsertRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"email"`
}
