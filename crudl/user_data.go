package crudl

import (
	"context"
	"database/sql"
	"fmt"
	"movie_backend_go/models"
	"strings"

	"github.com/google/uuid"
)

func CreateUserDataDB(ctx context.Context, db *sql.DB, userCreate models.CreateUserDataRequest) (models.UserData, error) {
	var createSchema = `
		INSERT INTO user_data(id, name, login, password) VALUES
		($1, $2, $3, $4)
		RETURNING id, name, login, password, created_at
		`
	res := db.QueryRow(createSchema, uuid.NewString(), userCreate.Name, userCreate.Login, userCreate.Password)
	if err := res.Err(); err != nil {
		return models.UserData{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	user := models.UserData{}
	err := res.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.UserData{}, fmt.Errorf("scanning created user: %w", err)
	}
	return user, nil
}

// Write correctly
func UpdateUserDataDB(ctx context.Context, db *sql.DB, userUpdate models.UpdateUserDataRequest, user_id string) (models.UserData, error) {
	var updateSchema = ` UPDATE user_data SET `
	updates := []string{}
	if userUpdate.Login != nil {
		updates = append(updates, fmt.Sprintf("login = '%s'", *userUpdate.Login))
	}
	if userUpdate.Name != nil {
		updates = append(updates, fmt.Sprintf("name = '%s'", *userUpdate.Name))
	}
	if userUpdate.Password != nil {
		updates = append(updates, fmt.Sprintf("password = '%s'", *userUpdate.Password))
	}
	updateString := strings.Join(updates, ", ")
	updateSchema += updateString
	updateSchema += fmt.Sprintf("\n WHERE id = '%s'", user_id)
	updateSchema += "\n RETURNING id, name, login, password, created_at"

	res := db.QueryRowContext(ctx, updateSchema)
	if err := res.Err(); err != nil {
		return models.UserData{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	user := models.UserData{}
	err := res.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.UserData{}, fmt.Errorf("scanning created user: %w", err)
	}
	return user, nil
}

func GetUserDB(ctx context.Context, db *sql.DB, id string) (models.UserData, error) {
	var getSchema = `
		SELECT id, name, login, password, created_at
		FROM user_data
		WHERE id = $1
		`
	res := db.QueryRow(getSchema, id)
	if err := res.Err(); err != nil {
		return models.UserData{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	user := models.UserData{}
	err := res.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.UserData{}, fmt.Errorf("scanning requested user: %w", err)
	}
	return user, nil
}

func GetUserListDB(ctx context.Context, db *sql.DB) (models.UserDataList, error) {
	var getMovieListSchema = `
		SELECT id, name, login, password, created_at
		FROM user_data
		`
	resRows, err := db.QueryContext(ctx, getMovieListSchema)
	if err != nil {
		return models.UserDataList{}, fmt.Errorf("get user list for user: %w", err)
	}
	defer resRows.Close()

	userList := models.UserDataList{}
	for resRows.Next() {
		select {
		case <-ctx.Done():
			return models.UserDataList{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			// Continue processing
		}

		var user = models.UserData{}
		if err := resRows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt); err != nil {
			return models.UserDataList{}, fmt.Errorf("scanning getting rows")
		}
		userList.UserList = append(userList.UserList, user)
	}
	if err := resRows.Err(); err != nil {
		return models.UserDataList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return userList, nil
}

func DeleteUserDB(ctx context.Context, db *sql.DB, id string) error {
	var deleteSchema = `
		DELETE FROM user_data
		WHERE id = $1
		`
	res, err := db.ExecContext(ctx, deleteSchema, id)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	return checkNonEmptyDeletion(res)

}
