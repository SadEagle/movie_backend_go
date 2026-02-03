package reqmodel

import "movie_backend_go/db/sqlc"

type MovieCreateRequest struct {
	Title string `json:"title"`
}

type MovieUpdateRequest struct {
	Title *string `json:"title"`
}
type MovieListResponse struct {
	MovieList []sqlc.Movie `json:"movie_list"`
}
