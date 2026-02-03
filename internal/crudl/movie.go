package crudl

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/db/sqlc"
)

func CreateMovie(ctx context.Context, querier sqlc.Querier, title string) (sqlc.Movie, error) {
	movie, err := querier.CreateMovie(ctx, title)
	return movie, err
}

func DeleteMovie(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) error {
	numDel, err := querier.DeleteMovie(ctx, movieID)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}

func GetMovie(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) (sqlc.GetMovieRow, error) {
	movie, err := querier.GetMovie(ctx, movieID)
	return movie, err
}

func GetMovieByTitle(ctx context.Context, querier sqlc.Querier, movieTitle string) (sqlc.GetMovieByTitleRow, error) {
	movie, err := querier.GetMovieByTitle(ctx, movieTitle)
	return movie, err
}

func GetMovieList(ctx context.Context, querier sqlc.Querier) ([]sqlc.Movie, error) {
	movieList, err := querier.GetMovieList(ctx)
	return movieList, err
}

func UpdateMovie(ctx context.Context, querier sqlc.Querier, movieUpdate sqlc.UpdateMovieParams) (sqlc.Movie, error) {
	movie, err := querier.UpdateMovie(ctx, movieUpdate)
	return movie, err
}
