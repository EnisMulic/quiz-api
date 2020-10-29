package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Quiz domain model
type Quiz struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name      string             `json:"name"`
	Timer     string             `json:"timer"`
	Questions []Question         `json:"questions"`
}

// QuizUpsertRequest model
type QuizUpsertRequest struct {
	Name   string             `json:"name"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Timer  string             `json:"timer"`
}
