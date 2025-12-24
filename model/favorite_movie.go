package model

import "github.com/google/uuid"

type FavoriteMovieBase struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type FavoriteMovie struct {
	ID uuid.UUID `json:"id"`
	MovieBase
}

type FavoriteMovieCreateRequest struct {
	MovieBase
}
type FavoriteMovieUpdateRequest struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type FavoriteMovieListResponse struct {
	favorite_movie []FavoriteMovie
}
