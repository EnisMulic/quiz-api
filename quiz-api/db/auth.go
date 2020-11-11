package db

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"log"
	"time"

	"github.com/EnisMulic/quiz-api/quiz-api/domain"
	"github.com/dgrijalva/jwt-go"
)

// Define the size of the salt
const saltSize = 16

// Generate a random 16 bytes securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// int the crypto.rand package
func generateRandomSalt(saltSize int) string {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return string(salt)
}

// Combine our password and salt and hash them using the SHA-512
// hashing algorithm and then return our hashed password
// as a base64 encoded string
func hashPassword(password string, salt string) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	var saltBytes = []byte(salt)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, saltBytes...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the sha512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a base64 encoded string
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

// Check if two passwords match
func doPasswordsMatch(passwordHash string, currPassword string, salt string) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return passwordHash == currPasswordHash
}

var secretKey = []byte("gosecretkey")

// GenerateJWT generate a json web token
func GenerateJWT(user domain.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = user.Email
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
