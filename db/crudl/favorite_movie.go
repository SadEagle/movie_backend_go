package crudl

import (
	"database/sql"
	"fmt"
	"movie_backend_go/models"
)

func CreateFavoriteMovieDB(db *sql.DB, favoriteMovieCreate models.CreateFavoriteMovieRequest) (models.FavoriteMovie, error) {
	var createShema = `
		INSERT INTO favorite_movie(user_id, movie_id)
		VALUES ($1, $2)
		RETURNING user_id, movie_id
		`
	res := db.QueryRow(createShema, favoriteMovieCreate.UserID, favoriteMovieCreate.MovieID)

	movie := models.FavoriteMovie{}
	err := res.Scan(&favoriteMovieCreate.UserID, &favoriteMovieCreate.MovieID)
	if err != nil {
		return models.FavoriteMovie{}, fmt.Errorf("scanning favorite_movie adding: %w", err)
	}
	return movie, nil
}

func GetFavoriteMovieListDB(db *sql.DB, user_id string) ([]string, error) {
	var getSchema = `
		SELECT array_agg(movie_id)
		FROM favorite_movie
		WHERE user_id = $1
		GROUP_BY user_id
		`
	res := db.QueryRow(getSchema, user_id)

	favoriteMovieIDList := []string{}
	err := res.Scan(&favoriteMovieIDList)
	if err != nil {
		return []string{}, fmt.Errorf("scanning user's favorite movie_id list: %w", err)
	}
	return favoriteMovieIDList, nil
}

func DeleteFavoriteMovieDB(db *sql.DB, user_id string, movie_id string) error {
	var deleteSchema = `
		DELETE FROM favorite_movie
		WHERE user_id = $1
		AND movie_id = $2
		`
	res, err := db.Exec(deleteSchema, user_id, movie_id)
	if err != nil {
		return fmt.Errorf("deleting movie from user's favorite: %w", err)
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
