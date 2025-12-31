package crudl

import (
	"database/sql"
	"fmt"
	"movie_backend_go/models"
	"strings"

	"github.com/google/uuid"
)

func CreateMovieDB(db *sql.DB, movieCreate models.CreateMovieRequest) (models.Movie, error) {
	var createShema = `
		INSERT INTO movie(id, title)
		VALUES ($1, $2)
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

func GetMovieListDB(db *sql.DB, movie_id string) (models.MovieListResponse, error) {
	var getListSchema = `
		SELECT movie_id
		FROM favorite_movie
		WHERE user_id = $1
		`
}

// Write correctly
func UpdateMovieDB(db *sql.DB, movieUpdate models.UpdateMovieRequest, movie_id string) (models.Movie, error) {
	var updateScheme = ` UPDATE movie SET `
	updates := []string{}
	if movieUpdate.Title != nil {
		updates = append(updates, fmt.Sprintf("title = '%s'", *movieUpdate.Title))
	}
	updateString := strings.Join(updates, ", ")
	updateScheme += updateString
	updateScheme += fmt.Sprintf("\n WHERE id = '%s'", movie_id)
	updateScheme += "\n RETURNING id, title, amount_marks, total_mark, rating, created_at"

	res := db.QueryRow(updateScheme)

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

	return checkNonEmptyDeletion(res)
}
