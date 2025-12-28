package models

import "github.com/google/uuid"

type MarkedMovie struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type CreateMarkedMovieRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type MarkedMovieResponse struct {
	UserID  uuid.UUID  `json:"user_id"`
	MovieID uuid.UUIDs `json:"movie_id"`
}

type MarkedMovieListResponse struct {
	UserID  uuid.UUID  `json:"user_id"`
	MovieID uuid.UUIDs `json:"movie_id"`
}
