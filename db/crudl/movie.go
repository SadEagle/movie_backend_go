package crudl

import (
	"database/sql"
	"fmt"
	"movie_backend_go/models"

	"github.com/google/uuid"
)

func CreateMovieDB(db *sql.DB, movieCreate models.CreateMovieRequest) (models.Movie, error) {
	var createShema = `
		INSERT INTO movie(id, title)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, amount_marks, total_mark, rating, created_at
		`
	res := db.QueryRow(createShema, uuid.NewString(), movieCreate.Title)

	movie := models.Movie{}
	err := res.Scan(&movie.ID, &movie.Title, &movie.AmountMarks, &movie.TotalMark, &movie.Rating, &movie.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("scanning created movie: %w", err)
	}
	return movie, nil
}

func GetMovieDB(db *sql.DB, id string) (models.Movie, error) {
	var getSchema = `
		SELECT id, title, amount_marks, total_mark, rating, created_at
		FROM movie
		WHERE id = $1
		`
	res := db.QueryRow(getSchema, id)

	movie := models.Movie{}
	err := res.Scan(&movie.ID, &movie.Title, &movie.AmountMarks, &movie.TotalMark, &movie.Rating, &movie.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("scanning requested movie: %w", err)
	}
	return movie, nil
}

func DeleteMovieDB(db *sql.DB, id string) error {
	var deleteSchema = `
		DELETE FROM movie
		WHERE id = $1
		`
	res, err := db.Exec(deleteSchema, id)
	if err != nil {
		return fmt.Errorf("deleting movie: %w", err)
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
