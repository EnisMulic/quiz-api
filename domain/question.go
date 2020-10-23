package domain

// Question domain model
type Question struct {
	ID        int      `json:"id"`
	Text      string   `json:"text"`
	Answers   []Answer `json:"answers"`
	CreatedOn string   `json:"-"`
	UpdatedOn string   `json:"-"`
	DeletedOn string   `json:"-"`
}
