package db

import (
	"context"
	"fmt"

	"github.com/EnisMulic/quiz-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// QuizRepository struct
type QuizRepository struct {
	collection *mongo.Collection
}

// NewQuizRepository ctor
func NewQuizRepository(c *mongo.Client) *QuizRepository {
	return &QuizRepository{c.Database("quiz-app").Collection("quiz")}
}

// Quizzes a collection of Quizzes
type Quizzes []*domain.Quiz

// GetQuizzes returns a slice of Quizzes
func (r *QuizRepository) GetQuizzes() Quizzes {
	var list Quizzes
	cur, err := r.collection.Find(nil, bson.M{})

	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}

	err = cur.All(context.Background(), &list)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}

	return list
}

// GetQuiz returns a single quiz
func (r *QuizRepository) GetQuiz(id primitive.ObjectID) domain.Quiz {
	var entity domain.Quiz
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return entity
}

// AddQuiz adds a new User
func (r *QuizRepository) AddQuiz(u *domain.Quiz) {
	_, err := r.collection.InsertOne(nil, u)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

// UpdateQuiz updates a user
func (r *QuizRepository) UpdateQuiz(id primitive.ObjectID, data map[string]interface{}) error {
	updateData := bson.M{
		"$set": data,
	}
	result, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, updateData)

	if result.MatchedCount != 1 {
		return ErrQuizNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteQuiz removes a quiz from the database
func (r *QuizRepository) DeleteQuiz(id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if result.DeletedCount != 1 {
		return ErrQuizNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// ErrQuizNotFound an error
var ErrQuizNotFound = fmt.Errorf("Quiz not found")
