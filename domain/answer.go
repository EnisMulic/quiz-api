package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Answer domain model
type Answer struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Text      string             `json:"text" validate:"required"`
	IsCorrect bool               `json:"is_correct" validate:"required"`
}
