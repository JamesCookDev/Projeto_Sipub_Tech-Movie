package domain

// Movie representa a entidade principal da nossa aplicação.
type Movie struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title"`
	Year  int    `json:"year" bson:"year"`
}