package crudl

import (
	"context"
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func CreateUser(ctx context.Context, querier sqlc.Querier, userCreate sqlc.CreateUserParams) (sqlc.UserDatum, error) {
	user, err := querier.CreateUser(ctx, userCreate)
	if err != nil {
		return sqlc.UserDatum{}, err
	}
	return user, nil
}

func DeleteUser(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) error {
	numDel, err := querier.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return EmptyDeletionError
	}
	return nil
}

func GetUserByID(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) (sqlc.UserDatum, error) {
	user, err := querier.GetUserByID(ctx, userID)
	if err != nil {
		return sqlc.UserDatum{}, err
	}
	return user, nil
}

func GetUserByLogin(ctx context.Context, querier sqlc.Querier, login string) (sqlc.UserDatum, error) {
	user, err := querier.GetUserByLogin(ctx, login)
	if err != nil {
		return sqlc.UserDatum{}, err
	}
	return user, nil
}

func GetUserList(ctx context.Context, querier sqlc.Querier) ([]sqlc.UserDatum, error) {
	userList, err := querier.GetUserList(ctx)
	if err != nil {
		return []sqlc.UserDatum{}, err
	}
	return userList, nil
}

func UpdateUser(ctx context.Context, querier sqlc.Querier, userUpdate sqlc.UpdateUserParams) (sqlc.UserDatum, error) {
	user, err := querier.UpdateUser(ctx, userUpdate)
	if err != nil {
		return sqlc.UserDatum{}, err
	}
	return user, nil
}
