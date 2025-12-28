package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	AmountMarks int       `json:"amount_marks"`
	TotalMark   int       `json:"total_mark"`
	Rating      float32   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

func (m *Movie) ToResponse() MovieResponse {
	return MovieResponse{
		ID:        m.ID,
		Title:     m.Title,
		Rating:    m.Rating,
		CreatedAt: m.CreatedAt,
	}
}

type CreateMovieRequest struct {
	Title string `json:"title"`
}

type UpdateMovieRequest struct {
	Title *string `json:"title"`
}

type MovieResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Rating    float32   `json:"rating"`
	CreatedAt time.Time `db:"created_at"`
}
