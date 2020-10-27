package handlers

import (
	"log"
	"net/http"

	"github.com/EnisMulic/quiz-api/db"
	"github.com/EnisMulic/quiz-api/domain"
)

// Auth struct
type Auth struct {
	l *log.Logger
	r *db.UserRepository
}

// NewAuth func
func NewAuth(l *log.Logger, r *db.UserRepository) *Auth {
	return &Auth{l, r}
}

// swagger:route POST /auth/register auth register
// Registers a user

// Register a new user
func (a *Auth) Register(rw http.ResponseWriter, r *http.Request) {
	user := new(domain.UserUpsertRequest)
	domain.FromJSON(user, r.Body)
	a.r.AddUser(user)

	err := domain.ToJSON(user, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

// swagger:route POST /auth/login auth login
// Login a user

// Login as a user
func (a *Auth) Login(rw http.ResponseWriter, r *http.Request) {
	user := new(domain.UserUpsertRequest)
	domain.FromJSON(user, r.Body)

	jwt, err := a.r.Login(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Write([]byte(`{"token":"` + jwt + `"}`))
}
