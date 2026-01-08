package crudl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movie_backend_go/models"
)

func AddFavoriteMovieDB(ctx context.Context, db *sql.DB, userID string, movieID string) (models.FavoriteMovie, error) {

	var createShema = `
		INSERT INTO favorite_movie(user_id, movie_id)
		VALUES ($1, $2)
		RETURNING user_id, movie_id
		`
	if userID == "" {
		return models.FavoriteMovie{}, errors.New("userID cannot be empty")
	}
	if movieID == "" {
		return models.FavoriteMovie{}, errors.New("userID cannot be empty")
	}

	res_create := db.QueryRowContext(ctx, createShema, userID, movieID)
	if err := res_create.Err(); err != nil {
		return models.FavoriteMovie{}, fmt.Errorf("check QueryRowContext correctness: %w", err)
	}

	favoriteMovie := models.FavoriteMovie{}
	err := res_create.Scan(&favoriteMovie.UserID, &favoriteMovie.MovieID)
	if err != nil {
		return models.FavoriteMovie{}, fmt.Errorf("scanning favorite_movie adding: %w", err)
	}
	return favoriteMovie, nil
}

func DeleteFavoriteMovieDB(ctx context.Context, db *sql.DB, userID string, movieID string) error {
	var createShema = `
		DELETE FROM favorite_movie
		WHERE user_id = $1 AND movie_id = $2
		`
	res, err := db.ExecContext(ctx, createShema, userID, movieID)
	if err != nil {
		return fmt.Errorf("delete favorite_movie: %w", err)
	}

	return checkNonEmptyDeletion(res)
}

// Simplified version. Better option will be create non-response datatype and later convert to response one
func GetFavoriteMovieListDB(ctx context.Context, db *sql.DB, userID string) (models.FavoriteMovieList, error) {
	var getListSchema = `
		SELECT user_id, movie_id
		FROM favorite_movie
		WHERE user_id = $1
		`

	resRows, err := db.QueryContext(ctx, getListSchema, userID)
	if err != nil {
		return models.FavoriteMovieList{}, fmt.Errorf("get favorite_movie list for user: %w", err)
	}
	defer resRows.Close()

	favMovieList := models.FavoriteMovieList{}
	for resRows.Next() {
		select {
		case <-ctx.Done():
			return models.FavoriteMovieList{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			// Continue processing
		}

		var favMovie models.FavoriteMovie
		if err := resRows.Scan(&favMovie.UserID, &favMovie.MovieID); err != nil {
			return models.FavoriteMovieList{}, fmt.Errorf("reading favorite movie list: %w", err)
		}
		favMovieList.FavMovieList = append(favMovieList.FavMovieList, favMovie)
	}
	if err := resRows.Err(); err != nil {
		return models.FavoriteMovieList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return favMovieList, err
}
