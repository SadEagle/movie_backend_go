package models

import "github.com/google/uuid"

type RatedMovie struct {
	UserID  uuid.UUID `json:"user_id"`
	MovieID uuid.UUID `json:"movie_id"`
	Rating  int       `json:"rating"`
}

type RatedMovieElem struct {
	MovieID uuid.UUID `json:"movie_id"`
	Rating  int       `json:"rating"`
}

type RatedMovieList struct {
	UserID         uuid.UUID        `json:"user_id"`
	RatedMovieList []RatedMovieElem `json:"rated_movie_list"`
}
