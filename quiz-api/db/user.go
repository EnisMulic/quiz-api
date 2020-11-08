package db

import (
	"context"
	"fmt"

	"github.com/EnisMulic/quiz-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository struct
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository ctor
func NewUserRepository(c *mongo.Client) *UserRepository {
	return &UserRepository{c.Database("quiz-app").Collection("users")}
}

// Users a collection of Users
type Users []*domain.User

// GetUsers returns a slice of Users
func (r *UserRepository) GetUsers() Users {
	var list Users
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

// GetUser returns a single user
func (r *UserRepository) GetUser(id primitive.ObjectID) domain.User {
	var user domain.User
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return user
}

// AddUser adds a new User
func (r *UserRepository) AddUser(u *domain.UserUpsertRequest) {
	var salt = generateRandomSalt(saltSize)
	var hash = hashPassword(u.Password, salt)

	user := map[string]string{
		"username": u.Username,
		"email":    u.Email,
		"salt":     salt,
		"hash":     hash,
	}

	_, err := r.collection.InsertOne(nil, user)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(id primitive.ObjectID, data *domain.UserUpsertRequest) error {
	user := new(domain.User)
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(user)
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%v", id)

	if data.Password != "" {
		if doPasswordsMatch(user.Hash, data.Password, user.Salt) == false {
			user.Salt = generateRandomSalt(saltSize)
			user.Hash = hashPassword(data.Password, user.Salt)
		}
	}

	updateData := bson.M{
		"username": data.Username,
		"email":    data.Email,
		"salt":     user.Salt,
		"hash":     user.Hash,
	}

	result, err := r.collection.ReplaceOne(context.TODO(), bson.M{"_id": id}, updateData)

	if result.MatchedCount != 1 {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser removes a user from the database
func (r *UserRepository) DeleteUser(id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})

	if result.DeletedCount != 1 {
		return ErrUserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

// Login return a jwt if password and email are correct
func (r *UserRepository) Login(request *domain.UserUpsertRequest) (string, error) {
	var dbUser domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": request.Email}).Decode(&dbUser)

	if doPasswordsMatch(dbUser.Hash, request.Password, dbUser.Salt) == false {
		return "", fmt.Errorf(`{"response":"Wrong Password!"}`)
	}

	jwtToken, err := GenerateJWT(dbUser)
	if err != nil {
		return "", fmt.Errorf(`{"message":"` + err.Error() + `"}`)
	}

	return jwtToken, nil
}

// ErrUserNotFound an error
var ErrUserNotFound = fmt.Errorf("User not found")
