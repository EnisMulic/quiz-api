package handlers

import (
	"log"
	"net/http"

	"main.go/domain"
	"main.go/repository"
)

// Users struct
type Users struct {
	l *log.Logger
}

// NewUser func
func NewUser(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		u.getUsers(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		u.addUser(rw, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getUsers returns the Users from the data store
func (u *Users) getUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")

	// fetch the products from the datastore
	list := repository.GetUsers()

	// serialize the list to JSON
	err := list.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (u *Users) addUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST Product")

	user := &domain.User{}

	err := user.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	repository.AddUser(user)

	u.l.Printf("User: %#v", user)
}
