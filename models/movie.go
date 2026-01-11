package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	AmountRates int       `json:"amount_rates"`
	Rating      float32   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

type MovieList struct {
	MovieList []Movie `json:"movie_list"`
}

type CreateMovieRequest struct {
	Title string `json:"title"`
}

type UpdateMovieRequest struct {
	Title *string `json:"title"`
}
