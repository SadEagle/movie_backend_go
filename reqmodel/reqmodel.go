package reqmodel

import (
	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/db/sqlc"
)

type FavoriteMovieListResponse struct {
	UserID           pgtype.UUID   `json:"user_id"`
	FavoriteMovieIDs []pgtype.UUID `json:"favorite_movie_ids"`
}

// Create and Update requests contain same possible parameters
type MovieRequest struct {
	Title string `json:"title"`
}

type MovieCommentRequest struct {
	UserID pgtype.UUID `json:"user_id"`
	Text   string      `json:"text"`
}

type UserCommentListResp struct {
	UserID          pgtype.UUID                  `json:"user_id"`
	UserCommentList []sqlc.GetUserCommentListRow `json:"user_comment_list"`
}

type MovieCommentListResp struct {
	MovieID          pgtype.UUID                   `json:"movie_id"`
	MovieCommentList []sqlc.GetMovieCommentListRow `json:"movie_comment_list"`
}

type MovieUpdateRequest struct {
	Title *string `json:"title"`
}
type MovieListResponse struct {
	MovieList []sqlc.Movie `json:"movie_list"`
}

type RatedMovieRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
	Rating  int16       `json:"rating"`
}

type RatedMovieUpdateRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
	Rating  *int16      `json:"rating"`
}

type RatedMovieListResponse struct {
	UserID         pgtype.UUID                  `json:"user_id"`
	RatedMovieList []sqlc.GetMovieRatingListRow `json:"rated_movie_list"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type UserUpdateRequest struct {
	Name     *string `json:"name"`
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type UserListResponse struct {
	UserList []sqlc.UserDatum `json:"user_list"`
}
