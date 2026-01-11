package crudl

import (
	"context"
	"database/sql"
	"fmt"
	"movie_backend_go/models"
	"strings"

	"github.com/google/uuid"
)

func baseGetUserDB(ctx context.Context, db *sql.DB, identField any, schema string) (models.User, error) {
	user := models.User{}
	resRow := db.QueryRowContext(ctx, schema, identField)
	err := resRow.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("scanning selected user data: %w", err)
	}
	return user, nil
}

func GetUserByIDDB(ctx context.Context, db *sql.DB, userID uuid.UUID) (models.User, error) {
	var userSearchSchema = `
		SELECT (id, name, login, password, is_admin created_at)
		FROM user_data
		WHERE id = $1
		`
	return baseGetUserDB(ctx, db, userID, userSearchSchema)
}

func GetUserByLoginDB(ctx context.Context, db *sql.DB, userLogin string) (models.User, error) {
	var userSearchSchema = `
		SELECT (id, name, login, password, is_admin, created_at)
		FROM user_data
		WHERE login = $1
		`
	return baseGetUserDB(ctx, db, userLogin, userSearchSchema)
}

func baseCreateUserDB(ctx context.Context, db *sql.DB, userCreate models.CreateUserRequest, is_admin bool) (models.User, error) {
	var createUserSchema = `
		INSERT INTO user_data(name, login, password, is_admin) VALUES
		($1, $2, $3, $4)
		RETURNING id, name, login, password, is_admin, created_at
		`
	res := db.QueryRow(createUserSchema, userCreate.Name, userCreate.Login, userCreate.Password, is_admin)

	user := models.User{}
	err := res.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("scanning created user: %w", err)
	}

	return user, nil
}

func CreateUserDB(ctx context.Context, db *sql.DB, userCreate models.CreateUserRequest) (models.User, error) {
	return baseCreateUserDB(ctx, db, userCreate, false)
}

func CreateAdminUserDB(ctx context.Context, db *sql.DB, userCreate models.CreateUserRequest) (models.User, error) {
	return baseCreateUserDB(ctx, db, userCreate, true)
}

// Write correctly
// FIX: SQL injection
func UpdateUserDB(ctx context.Context, db *sql.DB, userUpdate models.UpdateUserRequest, userID uuid.UUID) (models.User, error) {
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
	updateSchema += fmt.Sprintf("\n WHERE id = '%s'", userID)
	updateSchema += "\n RETURNING id, name, login, password, created_at"

	resRow := db.QueryRowContext(ctx, updateSchema)
	if err := resRow.Err(); err != nil {
		return models.User{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	user := models.User{}
	err := resRow.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("scanning created user: %w", err)
	}
	return user, nil
}

func GetUserListDB(ctx context.Context, db *sql.DB) (models.UserList, error) {
	var getMovieListSchema = `
		SELECT id, name, login, password, is_admin, created_at
		FROM user_data
		`
	resRows, err := db.QueryContext(ctx, getMovieListSchema)
	if err != nil {
		return models.UserList{}, fmt.Errorf("get user list for user: %w", err)
	}
	defer resRows.Close()

	userList := models.UserList{}
	for resRows.Next() {
		select {
		case <-ctx.Done():
			return models.UserList{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			// Continue processing
		}

		var user = models.User{}
		if err := resRows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt); err != nil {
			return models.UserList{}, fmt.Errorf("scanning getting rows")
		}
		userList.UserList = append(userList.UserList, user)
	}
	if err := resRows.Err(); err != nil {
		return models.UserList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return userList, nil
}

func DeleteUserDB(ctx context.Context, db *sql.DB, userID uuid.UUID) error {
	var deleteSchema = `
		DELETE FROM user_data
		WHERE id = $1
		`
	res, err := db.ExecContext(ctx, deleteSchema, userID)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	return checkNonEmptyDeletion(res)
}
