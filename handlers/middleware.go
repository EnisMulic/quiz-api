package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/EnisMulic/quiz-api/domain"
	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("gosecretkey")

// IsAuthorized checks if user is authorized
func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return secretKey, nil
			})

			if err != nil {
				fmt.Fprintf(rw, err.Error())
			}

			if token.Valid {
				next.ServeHTTP(rw, r)
			}
		} else {

			fmt.Fprintf(rw, "Not Authorized")
		}
	})
}

// KeyStruct a key
type KeyStruct struct{}

// MiddlewareValidateUser for validation
func (u Users) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := new(domain.UserUpsertRequest)
		err := domain.FromJSON(user, r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Error reading user", http.StatusBadRequest)
			return
		}

		// validate the user
		err = domain.Validate(user)
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
		ctx := context.WithValue(r.Context(), KeyStruct{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
