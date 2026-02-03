package reqmodel

import (
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type CommentCreateRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
	Text    string      `json:"text"`
}

type CommentUpdateRequest struct {
	Text string `json:"text"`
}

type MovieCommentRequest struct {
	MovieID pgtype.UUID `json:"movie_id"`
	Text    string      `json:"text"`
}

type UserCommentRequest struct {
	UserID pgtype.UUID `json:"user_id"`
	Text   string      `json:"text"`
}

type UserCommentListResponse struct {
	UserID          pgtype.UUID                  `json:"user_id"`
	UserCommentList []sqlc.GetUserCommentListRow `json:"user_comment_list"`
}

type MovieCommentListResponse struct {
	MovieID          pgtype.UUID                   `json:"movie_id"`
	MovieCommentList []sqlc.GetMovieCommentListRow `json:"movie_comment_list"`
}
