package crudl

import (
	"context"
	"database/sql"
	"movie_backend_go/models"

	"fmt"

	"github.com/google/uuid"
)

func GetRatedMovieListDB(ctx context.Context, db *sql.DB, userID uuid.UUID) (models.RatedMovieList, error) {
	ratedMovieListSchema := `
	SELECT movie_id, rating
	FROM rated_movie
	WHERE user_id = $1
	`

	resRows, err := db.QueryContext(ctx, ratedMovieListSchema, userID)
	if err != nil {
		return models.RatedMovieList{}, fmt.Errorf("get rated movie list for user: %w", err)
	}
	defer resRows.Close()

	ratedMovieList := models.RatedMovieList{UserID: userID}
	for resRows.Next() {
		select {
		case <-ctx.Done():
			return models.RatedMovieList{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			// Continue processing
		}

		ratedMovie := models.RatedMovieElem{}
		if err := resRows.Scan(&ratedMovie.MovieID); err != nil {
			return models.RatedMovieList{}, fmt.Errorf("reading rated movie list: %w", err)
		}
		ratedMovieList.RatedMovieList = append(ratedMovieList.RatedMovieList, ratedMovie)
	}
	if err := resRows.Err(); err != nil {
		return models.RatedMovieList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return ratedMovieList, err
}

func AddMovieRatingDB(ctx context.Context, db *sql.DB, userID uuid.UUID, movieData models.RatedMovieElem) (models.RatedMovie, error) {
	addRatedMovieSchema := `
	INSERT INTO rated_movie (user_id, movie_id, rating)
	VALUES ($1, $2, $3)
	RETURNING user_id, movie_id, rating
	`
	res_create := db.QueryRowContext(ctx, addRatedMovieSchema, userID, movieData.MovieID, movieData.Rating)
	if err := res_create.Err(); err != nil {
		return models.RatedMovie{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	ratedMovie := models.RatedMovie{}
	err := res_create.Scan(&ratedMovie.MovieID, &ratedMovie.Rating)
	if err != nil {
		return models.RatedMovie{}, fmt.Errorf("scanning added rated movie: %w", err)
	}
	return ratedMovie, nil
}

func DeleteRatedMovieDB(ctx context.Context, db *sql.DB, userID uuid.UUID, movieID uuid.UUID) error {
	var createShema = `
		DELETE FROM rated_movie
		WHERE user_id = $1 AND movie_id = $2
		`
	res, err := db.ExecContext(ctx, createShema, userID, movieID)
	if err != nil {
		return fmt.Errorf("delete rated movie: %w", err)
	}

	return checkNonEmptyDeletion(res)
}
