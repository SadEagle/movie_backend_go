package crudl

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/db/sqlc"
)

func GetUserFavoriteList(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) ([]pgtype.UUID, error) {
	favUserList, err := querier.GetUserFavoriteList(ctx, userID)
	return favUserList, err
}

func GetMovieFavoriteList(ctx context.Context, querier sqlc.Querier, movieID pgtype.UUID) ([]pgtype.UUID, error) {
	favMovieList, err := querier.GetUserFavoriteList(ctx, movieID)
	return favMovieList, err
}

func GetFavorite(ctx context.Context, querier sqlc.Querier, favGet sqlc.GetFavoriteParams) (sqlc.Favorite, error) {
	favorite, err := querier.GetFavorite(ctx, favGet)
	return favorite, err
}

func CreateFavorite(ctx context.Context, querier sqlc.Querier, favCreate sqlc.CreateFavoriteParams) (sqlc.Favorite, error) {
	favMovie, err := querier.CreateFavorite(ctx, favCreate)
	return favMovie, err
}

func DeleteFavorite(ctx context.Context, querier sqlc.Querier, favDelete sqlc.DeleteFavoriteParams) error {
	numDel, err := querier.DeleteFavorite(ctx, favDelete)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}
