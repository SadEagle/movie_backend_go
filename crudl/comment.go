package crudl

import (
	"context"
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func CreateComment(ctx context.Context, querier sqlc.Querier, commentCreate sqlc.CreateCommentParams) (sqlc.Comment, error) {
	movieComment, err := querier.CreateComment(ctx, commentCreate)
	return movieComment, err
}

func DeleteComment(ctx context.Context, querier sqlc.Querier, commentID pgtype.UUID) error {
	numDel, err := querier.DeleteComment(ctx, commentID)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}

func GetMovieCommentList(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) ([]sqlc.GetMovieCommentListRow, error) {
	movieCommentList, err := querier.GetMovieCommentList(ctx, movieID)
	return movieCommentList, err
}

func GetUserCommentList(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) ([]sqlc.GetUserCommentListRow, error) {
	userCommentList, err := querier.GetUserCommentList(ctx, userID)
	return userCommentList, err
}

func GetComment(ctx context.Context, querier sqlc.Querier, commentID pgtype.UUID) (sqlc.Comment, error) {
	comment, err := querier.GetComment(ctx, commentID)
	return comment, err
}

func UpdateMovieComment(ctx context.Context, querier sqlc.Querier, movieCommentUpdate sqlc.UpdateCommentParams) (sqlc.Comment, error) {
	movieComment, err := querier.UpdateComment(ctx, movieCommentUpdate)
	return movieComment, err
}
