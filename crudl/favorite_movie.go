package crudl

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/db/sqlc"
)

func CreateFavoriteMovie(ctx context.Context, querier sqlc.Querier, favMovieCreate sqlc.CreateFavoriteMovieParams) (sqlc.FavoriteMovie, error) {
	favMovie, err := querier.CreateFavoriteMovie(ctx, favMovieCreate)
	return favMovie, err
}

func DeleteFavoriteMovie(ctx context.Context, querier sqlc.Querier, favMovieDelete sqlc.DeleteFavoriteMovieParams) error {
	numDel, err := querier.DeleteFavoriteMovie(ctx, favMovieDelete)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}

func GetFavoriteMovieList(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) ([]pgtype.UUID, error) {
	favMovieIDList, err := querier.GetFavoriteMovieIDList(ctx, userID)
	return favMovieIDList, err
}
