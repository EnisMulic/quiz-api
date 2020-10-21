package domain

// Question domain model
type Question struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Answers   []Answer
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}
