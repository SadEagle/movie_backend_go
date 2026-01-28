package crudl

import (
	"context"
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func CreateMovieRating(ctx context.Context, querier sqlc.Querier, movieRatingCreate sqlc.CreateMovieRatingParams) (sqlc.RatedMovie, error) {
	movieRating, err := querier.CreateMovieRating(ctx, movieRatingCreate)
	return movieRating, err
}

func DeleteMovieRating(ctx context.Context, querier sqlc.Querier, movieRatingDelete sqlc.DeleteMovieRatingParams) error {
	numDel, err := querier.DeleteMovieRating(ctx, movieRatingDelete)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}

func GetMovieRatingList(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) ([]sqlc.GetMovieRatingListRow, error) {
	movieRatingList, err := querier.GetMovieRatingList(ctx, userID)
	return movieRatingList, err
}

func UpdateMovieRating(ctx context.Context, querier sqlc.Querier, movieRatingUpdate sqlc.UpdateMoveRatingParams) (sqlc.RatedMovie, error) {
	movieRating, err := querier.UpdateMoveRating(ctx, movieRatingUpdate)
	return movieRating, err
}
