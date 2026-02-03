package crudl

import (
	"context"
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetMovieRatingList(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) ([]sqlc.GetMovieRatingListRow, error) {
	movieRatingList, err := querier.GetMovieRatingList(ctx, movieID)
	return movieRatingList, err
}

func GetUserRatingList(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) ([]sqlc.GetUserRatingListRow, error) {
	userRatingList, err := querier.GetUserRatingList(ctx, userID)
	return userRatingList, err
}

func GetRating(ctx context.Context, querier sqlc.Querier, ratingGet sqlc.GetRatingParams) (sqlc.Rating, error) {
	rating, err := querier.GetRating(ctx, ratingGet)
	return rating, err
}

func CreateMovieRating(ctx context.Context, querier sqlc.Querier, movieRatingCreate sqlc.CreateRatingParams) (sqlc.Rating, error) {
	rating, err := querier.CreateRating(ctx, movieRatingCreate)
	return rating, err
}

func DeleteRating(ctx context.Context, querier sqlc.Querier, ratingDelete sqlc.DeleteRatingParams) error {
	numDel, err := querier.DeleteRating(ctx, ratingDelete)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}

func UpdateMovieRating(ctx context.Context, querier sqlc.Querier, ratingUpdate sqlc.UpdateRatingParams) (sqlc.Rating, error) {
	movieRating, err := querier.UpdateRating(ctx, ratingUpdate)
	return movieRating, err
}
