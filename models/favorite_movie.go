package models

import "github.com/google/uuid"

type FavoriteMovie struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type CreateFavoriteMovieRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type FavoriteMovieResponse struct {
	UserID  uuid.UUID  `json:"user_id"`
	MovieID uuid.UUIDs `json:"movie_id"`
}

type FavoriteMovieListResponse struct {
	UserID  uuid.UUID  `json:"user_id"`
	MovieID uuid.UUIDs `json:"movie_id"`
}
