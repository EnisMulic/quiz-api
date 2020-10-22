package handlers

import (
	"log"
	"net/http"
)

// Quizes struct
type Quizes struct {
	l *log.Logger
}

// NewQuiz a Quiz ctor
func NewQuiz(l *log.Logger) *Quizes {
	return &Quizes{l}
}

// swagger:route GET /quiz quiz listQuiz
// Returns a list of quizes

// GetQuizes HTTP Get Method
func (q *Quizes) GetQuizes(rw http.ResponseWriter, r *http.Request) {
	// ToBe Implemented
}

// swagger:route POST /quiz quiz createQuiz
// Create a new quiz

// AddQuiz adds a quiz to the "db"
func (q *Quizes) AddQuiz(rw http.ResponseWriter, r *http.Request) {
	// ToBe Implemented
}

// swagger:route PUT /quiz/{id} quiz updateQuiz
// Update a quiz

// UpdateQuiz updates a quiz
func (q *Quizes) UpdateQuiz(rw http.ResponseWriter, r *http.Request) {
	// ToBe Implemented
}
