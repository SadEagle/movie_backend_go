package crudl

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/db/sqlc"
)

func CreateMovieFavorite(ctx context.Context, querier sqlc.Querier, favMovieCreate sqlc.CreateMovieFavoriteParams) (sqlc.FavoriteMovie, error) {
	favMovie, err := querier.CreateMovieFavorite(ctx, favMovieCreate)
	if err != nil {
		return sqlc.FavoriteMovie{}, err
	}
	return favMovie, nil
}

func DeleteMovieFavorite(ctx context.Context, querier sqlc.Querier, favMovieDelete sqlc.DeleteMovieFavoriteParams) error {
	numDel, err := querier.DeleteMovieFavorite(ctx, favMovieDelete)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return EmptyDeletionError
	}
	return nil
}

func GetMovieFavoriteList(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) ([]pgtype.UUID, error) {
	favMovieList, err := querier.GetMovieFavoriteList(ctx, userID)
	if err != nil {
		return []pgtype.UUID{}, err
	}
	return favMovieList, nil
}
