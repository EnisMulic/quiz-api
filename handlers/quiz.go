package handlers

import (
	"log"
	"net/http"

	"github.com/EnisMulic/quiz-api/db"
	"github.com/EnisMulic/quiz-api/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuizzesResponse dto
//
// A list of users
// swagger:response QuizzesResponse
type QuizzesResponse struct {
	// All users in the db
	// in: body
	Body []domain.Quiz
}

// QuizResponse dto
//
// A user
// swagger:response QuizResponse
type QuizResponse struct {
	// A User in the db
	// in: body
	Body domain.Quiz
}

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
//
// responses:
//	200: QuizzesResponse

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
//
// responses:
//	200: QuizResponse

// GetQuiz returns the quiz from the data store
func (q *Quizes) GetQuiz(rw http.ResponseWriter, r *http.Request) {
	q.l.Println("Handle GET Quiz")

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

	entity := new(domain.QuizUpsertRequest)
	domain.FromJSON(entity, r.Body)
	q.r.AddQuiz(entity)

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

	user := new(domain.QuizUpsertRequest)
	domain.FromJSON(user, r.Body)

	err = q.r.UpdateQuiz(id, user)
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

// swagger:route POST /quiz/{id}/question quiz addQuestionToQuiz
// Create a new question in a quiz

// AddQuestion adds a question to a quiz
func (q *Quizes) AddQuestion(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	entity := new(domain.Question)
	domain.FromJSON(entity, r.Body)
	quiz, addErr := q.r.AddQuestion(id, *entity)

	if err != nil {
		http.Error(rw, addErr.Error(), http.StatusBadRequest)
		return
	}

	err = domain.ToJSON(quiz, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route DELETE /quiz/{id}/question/{question_id} quiz deleteQuestionFromQuiz
// Remove a question from a quiz

// DeleteQuestion deletes a question from a quiz
func (q *Quizes) DeleteQuestion(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	questionID, questionIDerr := primitive.ObjectIDFromHex(vars["question_id"])
	if questionIDerr != nil {
		http.Error(rw, "Unable to convert question_id", http.StatusBadRequest)
		return
	}

	quiz, deleteErr := q.r.DeleteQuestion(id, questionID)
	if deleteErr == db.ErrQuizNotFound {
		http.Error(rw, "Quiz not found", http.StatusNotFound)
		return
	}

	if deleteErr != nil {
		http.Error(rw, "Quiz not found", http.StatusInternalServerError)
		return
	}

	domain.ToJSON(quiz, rw)
}
