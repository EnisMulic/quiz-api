package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/EnisMulic/quiz-api/domain"
	"github.com/EnisMulic/quiz-api/repository"
	"github.com/gorilla/mux"
)

// A list of users
// swagger:response UsersResponse

// UsersResponse dto
type UsersResponse struct {
	// All users in the database
	// in: body
	Body []domain.User
}

// A user
// swagger:response UserResponse

// UserResponse dto
type UserResponse struct {
	// A User in the database
	// in: body
	Body domain.User
}

// Users struct
type Users struct {
	l *log.Logger
}

// NewUser func
func NewUser(l *log.Logger) *Users {
	return &Users{l}
}

// swagger:route GET /user user listUser
// Returns a list of users

// GetUsers returns the Users from the data store
func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")

	// fetch the products from the datastore
	list := repository.GetUsers()

	// serialize the list to JSON
	err := list.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route GET /user/{id} user singleUser
// Returns a single user

// GetUser returns the Users from the data store
func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET User")

	//To Be Implemented
}

// swagger:route POST /user user createUser
// Create a new user

// AddUser adds a user to the "db"
func (u *Users) AddUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST User")

	user := r.Context().Value(KeyUser{}).(domain.User)
	repository.AddUser(&user)

	u.l.Printf("User: %#v", user)
}

// swagger:route PUT /user/{id} user updateUser
// Update a users details

// UpdateUser updates a user
func (u Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	u.l.Println("Handle PUT User")
	user := r.Context().Value(KeyUser{}).(domain.User)

	err = repository.UpdateUser(id, &user)
	if err == repository.ErrUserNotFound {
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

// DeleteUser handles DELETE requests and removes items from the database
func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	//To Be Implemented
}

// KeyUser a key
type KeyUser struct{}

// MiddlewareValidateUser for validation
func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := domain.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the user
		err = user.Validate()
		if err != nil {
			u.l.Println("[ERROR] validating user", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating user: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
