package crudl

import (
	"context"
	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func IsLoginUserExist(ctx context.Context, querier sqlc.Querier, loginUser sqlc.IsLoginUserExistParams) (bool, error) {
	isExist, err := querier.IsLoginUserExist(ctx, loginUser)
	return isExist, err
}

func CreateUser(ctx context.Context, querier sqlc.Querier, userCreate sqlc.CreateUserParams) (sqlc.UserDatum, error) {
	user, err := querier.CreateUser(ctx, userCreate)
	return user, err
}

func DeleteUser(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) error {
	numDel, err := querier.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	if numDel == 0 {
		return ErrEmptyDeletion
	}
	return nil
}

func GetUserByID(ctx context.Context, querier sqlc.Querier, userID pgtype.UUID) (sqlc.UserDatum, error) {
	user, err := querier.GetUserByID(ctx, userID)
	return user, err
}

func GetUserByLogin(ctx context.Context, querier sqlc.Querier, login string) (sqlc.UserDatum, error) {
	user, err := querier.GetUserByLogin(ctx, login)
	return user, err
}

func GetUserList(ctx context.Context, querier sqlc.Querier) ([]sqlc.UserDatum, error) {
	userList, err := querier.GetUserList(ctx)
	return userList, err
}

func UpdateUser(ctx context.Context, querier sqlc.Querier, userUpdate sqlc.UpdateUserParams) (sqlc.UserDatum, error) {
	user, err := querier.UpdateUser(ctx, userUpdate)
	return user, err
}
