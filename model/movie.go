package model

import "github.com/google/uuid"

type MovieBase struct {
	Title  string `json:"title"`
	Rating string `json:"rating"`
}

type Movie struct {
	ID uuid.UUID `json:"id"`
	MovieBase
}

type MovieCreateRequest struct {
	MovieBase
}
type MovieUpdateRequest struct {
	Title  *string `json:"title"`
	Rating *string `json:"rating"`
}

type MovieResponse struct {
	Movie
}
