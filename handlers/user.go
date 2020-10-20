// Package classification of Quiz API
//
// Documentation for Quiz API
//
//  Schames: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//	swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"main.go/domain"
	"main.go/repository"
)

// // A list of users
// // swagger:response usersResponse
// type usersResponse struct {
// 	// All users in the database
// 	// in: Body
// 	Body []domain.User
// }

// Users struct
type Users struct {
	l *log.Logger
}

// NewUser func
func NewUser(l *log.Logger) *Users {
	return &Users{l}
}

// swagger:route GET /user users listUser
// Returns a list of user

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

// swagger:route POST /user users createUser
// Create a new user

// AddUser adds a user to the "db"
func (u *Users) AddUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST Product")

	user := r.Context().Value(KeyUser{}).(domain.User)
	repository.AddUser(&user)

	u.l.Printf("User: %#v", user)
}

// swagger:route PUT /user/{id} users updateUser
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
