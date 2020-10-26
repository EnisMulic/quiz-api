package db

import (
	"context"
	"fmt"

	"github.com/EnisMulic/quiz-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection string = "users"

// UserRepository struct
type UserRepository struct {
	c *mongo.Client
}

// NewUserRepository ctor
func NewUserRepository(c *mongo.Client) *UserRepository {
	return &UserRepository{c}
}

// Users a collection of Users
type Users []*domain.User

// GetUsers returns a slice of Users
func (ur *UserRepository) GetUsers() Users {
	collection := ur.c.Database("quiz-app").Collection(userCollection)

	var list Users
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

// GetUser returns a single user
func (ur *UserRepository) GetUser(id primitive.ObjectID) domain.User {
	collection := ur.c.Database("quiz-app").Collection(userCollection)

	var user domain.User
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return user
}

// AddUser adds a new User
func (ur *UserRepository) AddUser(u *domain.UserUpsertRequest) {
	_, err := ur.c.Database("quiz-app").Collection(userCollection).InsertOne(nil, u)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

// UpdateUser updates a user
func (ur *UserRepository) UpdateUser(id primitive.ObjectID, data *domain.UserUpsertRequest) error {
	collection := ur.c.Database("quiz-app").Collection(userCollection)

	updateData := bson.D{
		{"$set", data.Username},
	}
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, updateData)

	if result.MatchedCount != 1 {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser removes a user from the database
func (ur *UserRepository) DeleteUser(id primitive.ObjectID) error {
	collection := ur.c.Database("quiz-app").Collection(userCollection)

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if result.DeletedCount != 1 {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// ErrUserNotFound an error
var ErrUserNotFound = fmt.Errorf("User not found")
