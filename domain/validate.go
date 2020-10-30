package domain

import "github.com/go-playground/validator"

// Validate struct
func Validate(i interface{}) error {
	validate := validator.New()
	return validate.Struct(i)
}
