package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EnisMulic/quiz-api/db"
	"github.com/EnisMulic/quiz-api/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KeyQuiz a key
type KeyQuiz struct{}

// Quizes struct
type Quizes struct {
	l *log.Logger
	r *db.QuizRepository
}

// NewQuiz a Quiz ctor
func NewQuiz(l *log.Logger, r *db.QuizRepository) *Quizes {
	return &Quizes{l, r}
}

// swagger:route GET /quiz quiz listQuiz
// Returns a list of quizes

// GetQuizes HTTP Get Method
func (q *Quizes) GetQuizes(rw http.ResponseWriter, r *http.Request) {
	q.l.Println("Handle GET Quizzes")

	// fetch the products from the datastore

	list := q.r.GetQuizzes()

	// serialize the list to JSON
	err := domain.ToJSON(list, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route GET /quiz/{id} quiz singleQuiz
// Returns a single quiz

// GetQuiz returns the quiz from the data store
func (q *Quizes) GetQuiz(rw http.ResponseWriter, r *http.Request) {
	q.l.Println("Handle GET User")

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	entity := q.r.GetQuiz(id)

	// serialize the list to JSON
	err = domain.ToJSON(entity, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route POST /quiz quiz createQuiz
// Create a new quiz

// AddQuiz adds a quiz to the "db"
func (q *Quizes) AddQuiz(rw http.ResponseWriter, r *http.Request) {
	q.l.Println("Handle POST Quiz")

	entity := r.Context().Value(KeyQuiz{}).(domain.Quiz)
	q.r.AddQuiz(&entity)

	q.l.Printf("User: %#v", entity)
}

// swagger:route PUT /quiz/{id} quiz updateQuiz
// Update a quiz

// UpdateQuiz updates a quiz
func (q *Quizes) UpdateQuiz(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(rw, "json body is incorrect", http.StatusBadRequest)
		return
	}

	err = q.r.UpdateQuiz(id, updateData)
	if err == db.ErrQuizNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /quiz/{id} quiz deleteQuiz
// Delete a quiz

// DeleteQuiz handles DELETE requests and removes items from the database
func (q *Quizes) DeleteQuiz(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	err = q.r.DeleteQuiz(id)
	if err == db.ErrQuizNotFound {
		http.Error(rw, "Quiz not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Quiz not found", http.StatusInternalServerError)
		return
	}
}
