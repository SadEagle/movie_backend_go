package crudl

import (
	"database/sql"
	"fmt"
	"movie_backend_go/models"

	"github.com/google/uuid"
)

func CreateUserDataDB(db *sql.DB, userCreate models.CreateUserDataRequest) (models.UserData, error) {
	var createSchema = `
		INSERT INTO user_data(id, name, login, password) VALUES
		($1, $2, $3, $4)
		RETURNING id, name, login, password, created_at
		`
	res := db.QueryRow(createSchema, uuid.NewString(), userCreate.Name, userCreate.Login, userCreate.Password)

	user := models.UserData{}
	err := res.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.UserData{}, fmt.Errorf("scanning created user: %w", err)
	}
	return user, nil
}

func GetUserDB(db *sql.DB, id string) (models.UserData, error) {
	var getSchema = `
		SELECT id, name, login, password, created_at
		FROM user_data
		WHERE id = $1
		`
	res := db.QueryRow(getSchema, id)

	user := models.UserData{}
	err := res.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.UserData{}, fmt.Errorf("scanning requested user: %w", err)
	}
	return user, nil
}

func DeleteUserDB(db *sql.DB, id string) error {
	var deleteSchema = `
		DELETE FROM user_data
		WHERE id = $1
		`
	res, err := db.Exec(deleteSchema, id)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("calculate affected rows by delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("0 rows were deleted")
	}
	return nil
}
