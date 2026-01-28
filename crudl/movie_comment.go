package crudl

import (
	"context"
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func CreateMovieComment(ctx context.Context, querier sqlc.Querier, movieCommentCreate sqlc.CreateMovieCommentParams) (sqlc.MovieComment, error) {
	movieComment, err := querier.CreateMovieComment(ctx, movieCommentCreate)
	return movieComment, err
}

func DeleteMovieComment(ctx context.Context, querier sqlc.Querier, movieCommentDelete sqlc.DeleteMovieCommentParams) error {
	numDel, err := querier.DeleteMovieComment(ctx, movieCommentDelete)
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

func UpdateMovieComment(ctx context.Context, querier sqlc.Querier, movieCommentUpdate sqlc.UpdateMovieCommentParams) (sqlc.MovieComment, error) {
	movieComment, err := querier.UpdateMovieComment(ctx, movieCommentUpdate)
	return movieComment, err
}
