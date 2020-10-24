package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Answer domain model
type Answer struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Text       string             `json:"text"`
	IsCorrect  bool               `json:"is_correct"`
	QuestionID primitive.ObjectID `json:"question_id"`
}
