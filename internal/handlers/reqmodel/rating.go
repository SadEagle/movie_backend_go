package reqmodel

import (
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type RatingCreateRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
	Rating  int16       `json:"rating"`
}

type RatingUpdateRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
	Rating  *int16      `json:"rating"`
}

type RatingMyDeleteRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
}

type RatingDeleteRequest struct {
	UserID  pgtype.UUID `json:"user_id"`
	MovieID pgtype.UUID `json:"movie_id"`
}

type UserRatingListResponse struct {
	UserID         pgtype.UUID                 `json:"user_id"`
	UserRatingList []sqlc.GetUserRatingListRow `json:"user_rating_list"`
}

type MovieRatingListResponse struct {
	MovieID         pgtype.UUID                  `json:"movie_id"`
	MovieRatingList []sqlc.GetMovieRatingListRow `json:"movie_rating_list"`
}
