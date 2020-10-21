package domain

// Answer domain model
type Answer struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	QuestionID int    `json:"question_id"`
	CreatedOn  string `json:"-"`
	UpdatedOn  string `json:"-"`
	DeletedOn  string `json:"-"`
}
