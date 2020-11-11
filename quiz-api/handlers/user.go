package handlers

import (
	"log"
	"net/http"

	"github.com/EnisMulic/quiz-api/quiz-api/db"
	"github.com/EnisMulic/quiz-api/quiz-api/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UsersResponse dto
//
// A list of users
// swagger:response UsersResponse
type UsersResponse struct {
	// All users in the db
	// in: body
	Body []domain.User
}

// UserResponse dto
//
// A user
// swagger:response UserResponse
type UserResponse struct {
	// A User in the db
	// in: body
	Body domain.User
}

// Users struct
type Users struct {
	l *log.Logger
	r *db.UserRepository
}

// NewUser func
func NewUser(l *log.Logger, r *db.UserRepository) *Users {
	return &Users{l, r}
}

// swagger:route GET /user user listUser
// Returns a list of users
//
// responses:
//	200: UsersResponse

// GetUsers returns the Users from the data store
func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")

	// fetch the products from the datastore

	list := u.r.GetUsers()

	// serialize the list to JSON
	err := domain.ToJSON(list, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route GET /user/{id} user singleUser
// Returns a single user
//
// responses:
//	200: UserResponse

// GetUser returns the Users from the data store
func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET User")

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	user := u.r.GetUser(id)

	// serialize the list to JSON
	err = domain.ToJSON(user, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route POST /user user createUser
// Create a new user

// AddUser adds a user to the "db"
func (u *Users) AddUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST User")

	user := new(domain.UserUpsertRequest)
	domain.FromJSON(user, r.Body)
	u.r.AddUser(user)

	err := domain.ToJSON(user, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route PUT /user/{id} user updateUser
// Update a users details
//
// responses:
//	404: noContentResponse

// UpdateUser updates a user
func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle PUT User")

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	user := new(domain.UserUpsertRequest)
	domain.FromJSON(user, r.Body)

	err = u.r.UpdateUser(id, user)
	if err == db.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /user/{id} user deleteUser
// Delete a user
//
// responses:
//	404: noContentResponse

// DeleteUser handles DELETE requests and removes items from the db
func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	err = u.r.DeleteUser(id)
	if err == db.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}
