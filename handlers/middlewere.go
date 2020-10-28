package handlers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("gosecretkey")

// IsAuthorized checks if user is authorized
// func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		if r.Header["Token"] != nil {

// 			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("There was an error")
// 				}
// 				return secretKey, nil
// 			})

// 			if err != nil {
// 				fmt.Fprintf(w, err.Error())
// 			}

// 			if token.Valid {
// 				endpoint(w, r)
// 			}
// 		} else {

// 			fmt.Fprintf(w, "Not Authorized")
// 		}
// 	})
// }

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
