package db

import (
	"context"
	"fmt"

	"github.com/EnisMulic/quiz-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var quizCollection string = "quiz"

// QuizRepository struct
type QuizRepository struct {
	c *mongo.Client
}

// NewQuizRepository ctor
func NewQuizRepository(c *mongo.Client) *QuizRepository {
	return &QuizRepository{c}
}

// Quizzes a collection of Quizzes
type Quizzes []*domain.Quiz

// GetQuizzes returns a slice of Quizzes
func (ur *QuizRepository) GetQuizzes() Quizzes {
	collection := ur.c.Database("quiz-app").Collection(quizCollection)

	var list Quizzes
	cur, err := collection.Find(nil, bson.M{})

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
func (ur *QuizRepository) GetQuiz(id primitive.ObjectID) domain.Quiz {
	collection := ur.c.Database("quiz-app").Collection(quizCollection)

	var entity domain.Quiz
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return entity
}

// AddQuiz adds a new User
func (ur *QuizRepository) AddQuiz(u *domain.Quiz) {
	_, err := ur.c.Database("quiz-app").Collection(quizCollection).InsertOne(nil, u)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

// UpdateQuiz updates a user
func (ur *QuizRepository) UpdateQuiz(id primitive.ObjectID, data map[string]interface{}) error {
	collection := ur.c.Database("quiz-app").Collection(quizCollection)

	updateData := bson.M{
		"$set": data,
	}
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, updateData)

	if result.MatchedCount != 1 {
		return ErrQuizNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteQuiz removes a quiz from the database
func (ur *QuizRepository) DeleteQuiz(id primitive.ObjectID) error {
	collection := ur.c.Database("quiz-app").Collection(quizCollection)

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

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
