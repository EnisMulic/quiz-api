package domain

// Quiz domain model
type Quiz struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Timer     string     `json:"timer"`
	Questions []Question `json:"questions"`
	CreatedOn string     `json:"-"`
	UpdatedOn string     `json:"-"`
	DeletedOn string     `json:"-"`
}
