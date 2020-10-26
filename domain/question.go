package domain

// QuestionType Enumeration
type QuestionType int

// QuestionType Enumeration Implementation
const (
	TrueFalse QuestionType = iota
	SingleAnswer
	MultipleAnswer
)

// Question domain model
type Question struct {
	Text         string       `json:"text"`
	Answers      []Answer     `json:"answers"`
	QuestionType QuestionType `json:"type"`
}
