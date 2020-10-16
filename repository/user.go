package repository

import (
	"encoding/json"
	"io"
	"time"

	"main.go/domain"
)

// example data source
var userList = []*domain.User{
	&domain.User{
		ID:        1,
		Username:  "User1",
		Email:     "user1@email.com",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	&domain.User{
		ID:        2,
		Username:  "User2",
		Email:     "user2@email.com",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}

// Users a collection of Users
type Users []*domain.User

// ToJSON serializes the contents of the collection to JSON
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// GetUsers returns a slice of Users
func GetUsers() Users {
	return userList
}

// AddUser adds a new User
func AddUser(u *domain.User) {
	userList = append(userList, u)
}
