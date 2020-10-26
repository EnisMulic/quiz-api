package handlers

import (
	"log"
	"net/http"

	"github.com/EnisMulic/quiz-api/db"
	"github.com/EnisMulic/quiz-api/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A list of users
// swagger:response UsersResponse

// UsersResponse dto
type UsersResponse struct {
	// All users in the db
	// in: body
	Body []domain.User
}

// A user
// swagger:response UserResponse

// UserResponse dto
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

// UpdateUser updates a user
func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle PUT User")

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	entity := new(domain.UserUpsertRequest)
	domain.FromJSON(entity, r.Body)

	err = u.r.UpdateUser(id, entity)
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

// KeyUser a key
type KeyUser struct{}

// MiddlewareValidateUser for validation
// func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		user := domain.User{}
// 		fmt.Printf("%v", r.Body)
// 		err := domain.FromJSON(user, r.Body)
// 		if err != nil {
// 			u.l.Println("[ERROR] deserializing user", err)
// 			http.Error(rw, "Error reading user", http.StatusBadRequest)
// 			return
// 		}

// 		// validate the user
// 		err = user.Validate()
// 		if err != nil {
// 			u.l.Println("[ERROR] validating user", err)
// 			http.Error(
// 				rw,
// 				fmt.Sprintf("Error validating user: %s", err),
// 				http.StatusBadRequest,
// 			)
// 			return
// 		}

// 		// add the product to the context
// 		ctx := context.WithValue(r.Context(), KeyUser{}, user)
// 		r = r.WithContext(ctx)

// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(rw, r)
// 	})
// }
