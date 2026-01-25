package crudl

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/db/sqlc"
)

func CreateMovie(ctx context.Context, querier sqlc.Querier, title string) (sqlc.Movie, error) {
	movie, err := querier.CreateMovie(ctx, title)
	if err != nil {
		return sqlc.Movie{}, err
	}
	return movie, err
}

func DeleteMovie(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) error {
	numDel, err := querier.DeleteMovie(ctx, movieID)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return EmptyDeletionError
	}
	return nil
}

func GetMovieByID(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) (sqlc.GetMovieByIDRow, error) {
	movie, err := querier.GetMovieByID(ctx, movieID)
	if err != nil {
		return sqlc.GetMovieByIDRow{}, err
	}
	return movie, nil
}

func GetMovieByTitle(ctx context.Context, querier sqlc.Querier, movieTitle string) (sqlc.GetMovieByTitleRow, error) {
	movie, err := querier.GetMovieByTitle(ctx, movieTitle)
	if err != nil {
		return sqlc.GetMovieByTitleRow{}, err
	}
	return movie, nil
}

func GetMovieList(ctx context.Context, querier sqlc.Querier) ([]sqlc.Movie, error) {
	movieList, err := querier.GetMovieList(ctx)
	if err != nil {
		return []sqlc.Movie{}, err
	}
	return movieList, nil
}

func UpdateMovie(ctx context.Context, querier sqlc.Querier, movieUpdate sqlc.UpdateMovieParams) (sqlc.Movie, error) {
	movie, err := querier.UpdateMovie(ctx, movieUpdate)
	if err != nil {
		return sqlc.Movie{}, err
	}
	return movie, nil
}
