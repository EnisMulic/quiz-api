package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// QuestionType Enumeration
type QuestionType int

// QuestionType Enumeration Implementation
const (
	TrueFalse QuestionType = iota
	SingleAnswer
	MultipleAnswer
)

// Question domain model
type Question struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Text         string             `json:"text" validate:"required"`
	Answers      []Answer           `json:"answers" validate:"required"`
	QuestionType QuestionType       `json:"type" validate:"required"`
}
