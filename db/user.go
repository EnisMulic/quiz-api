package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/EnisMulic/quiz-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// example data source
var userList = []*domain.User{
	&domain.User{
		Username:  "User1",
		Email:     "user1@email.com",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	&domain.User{
		Username:  "User2",
		Email:     "user2@email.com",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}

var userCollection string = "users"

// UserRepository struct
type UserRepository struct {
	c *mongo.Client
}

// NewUserRepository ctor
func NewUserRepository(c *mongo.Client) *UserRepository {
	return &UserRepository{c}
}

// ToJSON serializes the contents of the collection to JSON
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
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
func (ur *UserRepository) AddUser(u *domain.User) domain.User {
	result, err := ur.c.Database("quiz-app").Collection(userCollection).InsertOne(nil, u)
	if err != nil {
		fmt.Printf("%s", err)
	}

	u.ID = result.InsertedID.(primitive.ObjectID)
	return *u
}

// UpdateUser updates a user
func (ur *UserRepository) UpdateUser(id primitive.ObjectID, data map[string]interface{}) error {
	collection := ur.c.Database("quiz-app").Collection(userCollection)

	updateData := bson.M{
		"$set": data,
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

// ErrUserNotFound an error
var ErrUserNotFound = fmt.Errorf("User not found")
