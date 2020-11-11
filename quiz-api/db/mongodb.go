package db

import (
	"github.com/EnisMulic/quiz-api/quiz-api/domain"
	"github.com/goonode/mogo"
)

// RegisterModels mongo db model registry
func RegisterModels() {
	mogo.ModelRegistry.Register(domain.User{})
	mogo.ModelRegistry.Register(domain.Quiz{})
	mogo.ModelRegistry.Register(domain.Question{})
	mogo.ModelRegistry.Register(domain.Answer{})
}
