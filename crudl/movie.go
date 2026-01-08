package crudl

import (
	"context"
	"database/sql"
	"fmt"
	"movie_backend_go/models"
	"strings"

	"github.com/google/uuid"
)

func CreateMovieDB(ctx context.Context, db *sql.DB, movieCreate models.CreateMovieRequest) (models.Movie, error) {
	var createShema = `
		INSERT INTO movie(id, title)
		VALUES ($1, $2)
		RETURNING id, title, rating, created_at
		`
	res := db.QueryRowContext(ctx, createShema, uuid.NewString(), movieCreate.Title)
	if err := res.Err(); err != nil {
		return models.Movie{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	movie := models.Movie{}
	err := res.Scan(&movie.ID, &movie.Title, &movie.Rating, &movie.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("scanning created movie: %w", err)
	}
	return movie, nil
}

// Write correctly
func UpdateMovieDB(ctx context.Context, db *sql.DB, movieUpdate models.UpdateMovieRequest, movieID string) (models.Movie, error) {
	var updateScheme = ` UPDATE movie SET `
	updates := []string{}
	if movieUpdate.Title != nil {
		updates = append(updates, fmt.Sprintf("title = '%s'", *movieUpdate.Title))
	}
	updateString := strings.Join(updates, ", ")
	updateScheme += updateString
	updateScheme += fmt.Sprintf("\n WHERE id = '%s'", movieID)
	updateScheme += "\n RETURNING id, title, rating, created_at"

	res := db.QueryRowContext(ctx, updateScheme)
	if err := res.Err(); err != nil {
		return models.Movie{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	movie := models.Movie{}
	err := res.Scan(&movie.ID, &movie.Title, &movie.Rating, &movie.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("scanning created movie: %w", err)
	}
	return movie, nil
}

func GetMovieDB(ctx context.Context, db *sql.DB, movieID string) (models.Movie, error) {
	var getSchema = `
		SELECT id, title, rating, created_at
		FROM movie
		WHERE id = $1
		`
	res := db.QueryRowContext(ctx, getSchema, movieID)
	if err := res.Err(); err != nil {
		return models.Movie{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	movie := models.Movie{}
	err := res.Scan(&movie.ID, &movie.Title, &movie.Rating, &movie.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("scanning requested movie: %w", err)
	}
	return movie, nil
}

func GetMovieListDB(ctx context.Context, db *sql.DB) (models.MovieList, error) {
	var getMovieListSchema = `
		SELECT id, title, rating, created_at
		FROM movie
		`
	resRows, err := db.QueryContext(ctx, getMovieListSchema)
	if err != nil {
		return models.MovieList{}, fmt.Errorf("get movie list for user: %w", err)
	}
	defer resRows.Close()

	movieList := models.MovieList{}
	for resRows.Next() {
		select {
		case <-ctx.Done():
			return models.MovieList{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			// Continue processing
		}

		var movie = models.Movie{}
		if err := resRows.Scan(&movie.ID, &movie.Title, &movie.Rating, &movie.CreatedAt); err != nil {
			return models.MovieList{}, fmt.Errorf("scanning getting rows")
		}
		movieList.MovieList = append(movieList.MovieList, movie)
	}
	if err := resRows.Err(); err != nil {
		return models.MovieList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return movieList, nil
}

func DeleteMovieDB(ctx context.Context, db *sql.DB, movieID string) error {
	var deleteSchema = `
		DELETE FROM movie
		WHERE id = $1
		`
	res, err := db.ExecContext(ctx, deleteSchema, movieID)
	if err != nil {
		return fmt.Errorf("deleting movie: %w", err)
	}

	return checkNonEmptyDeletion(res)
}
