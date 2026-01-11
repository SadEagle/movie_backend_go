package models

import "github.com/google/uuid"

type FavoriteMovie struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
}

type FavoriteMovieElem struct {
	MovieID uuid.UUID `json:"movie_id"`
}

type FavoriteMovieList struct {
	UserID       uuid.UUID           `json:"user_id"`
	FavMovieList []FavoriteMovieElem `json:"favorite_movie_list"`
}
