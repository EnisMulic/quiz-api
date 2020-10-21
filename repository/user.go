package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/EnisMulic/quiz-api/domain"
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

// UpdateUser updates a user
func UpdateUser(id int, u *domain.User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	u.ID = id
	userList[pos] = u

	return nil
}

// ErrUserNotFound an error
var ErrUserNotFound = fmt.Errorf("User not found")

func findUser(id int) (*domain.User, int, error) {
	for i, p := range userList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrUserNotFound
}

func getNextID() int {
	lp := userList[len(userList)-1]
	return lp.ID + 1
}
