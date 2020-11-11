package db

import (
	"context"
	"fmt"

	"github.com/EnisMulic/quiz-api/quiz-api/domain"
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
func (r *QuizRepository) AddQuiz(u *domain.QuizUpsertRequest) {
	_, err := r.collection.InsertOne(nil, u)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

// UpdateQuiz updates a user
func (r *QuizRepository) UpdateQuiz(id primitive.ObjectID, data *domain.QuizUpsertRequest) error {
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

// AddQuestion to quiz
func (r *QuizRepository) AddQuestion(id primitive.ObjectID, question domain.Question) (domain.Quiz, error) {
	question.ID = primitive.NewObjectID()
	var entity domain.Quiz
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		return domain.Quiz{}, err
	}

	entity.Questions = append(entity.Questions, question)
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": entity})

	if err != nil {
		return domain.Quiz{}, err
	}

	return entity, nil
}

// DeleteQuestion from quiz
func (r *QuizRepository) DeleteQuestion(id primitive.ObjectID, questionID primitive.ObjectID) (domain.Quiz, error) {
	var entity domain.Quiz
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		return domain.Quiz{}, err
	}

	for i, v := range entity.Questions {
		if v.ID == questionID {
			entity.Questions = append(entity.Questions[:i], entity.Questions[i+1:]...)
			break
		}
	}

	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": entity})
	if err != nil {
		return domain.Quiz{}, err
	}

	return entity, nil
}

// ErrQuizNotFound an error
var ErrQuizNotFound = fmt.Errorf("Quiz not found")
