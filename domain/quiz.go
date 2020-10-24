package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Quiz domain model
type Quiz struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	// User      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name      string     `json:"name"`
	Timer     string     `json:"timer"`
	Questions []Question `json:"questions"`
}
