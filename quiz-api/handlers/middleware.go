package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/EnisMulic/quiz-api/quiz-api/domain"
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

// MiddlewareValidateQuiz for validation
func (q Quizes) MiddlewareValidateQuiz(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		quiz := new(domain.QuizUpsertRequest)
		err := domain.FromJSON(quiz, r.Body)
		if err != nil {
			q.l.Println("[ERROR] deserializing quiz", err)
			http.Error(rw, "Error reading quiz", http.StatusBadRequest)
			return
		}

		// validate the quiz
		err = domain.Validate(quiz)
		if err != nil {
			q.l.Println("[ERROR] validating quiz", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating quiz: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyStruct{}, quiz)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

// MiddlewareValidateQuestion for validation
func (q Quizes) MiddlewareValidateQuestion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		question := new(domain.Question)
		err := domain.FromJSON(question, r.Body)
		if err != nil {
			q.l.Println("[ERROR] deserializing question", err)
			http.Error(rw, "Error reading question", http.StatusBadRequest)
			return
		}

		// validate the question
		err = domain.Validate(question)
		if err != nil {
			q.l.Println("[ERROR] validating question", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating question: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyStruct{}, question)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
