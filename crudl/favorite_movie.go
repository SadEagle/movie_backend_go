package crudl

import (
	"context"
	"database/sql"
	"fmt"
	"movie_backend_go/models"

	"github.com/google/uuid"
)

func AddFavoriteMovieDB(ctx context.Context, db *sql.DB, userID uuid.UUID, movieID uuid.UUID) (models.FavoriteMovie, error) {

	var addFavoriteMovieSchema = `
		INSERT INTO favorite_movie(user_id, movie_id)
		VALUES ($1, $2)
		RETURNING user_id, movie_id
		`

	res_create := db.QueryRowContext(ctx, addFavoriteMovieSchema, userID, movieID)
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

func DeleteFavoriteMovieDB(ctx context.Context, db *sql.DB, userID uuid.UUID, movieID uuid.UUID) error {
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

// TODO: Probably, optimize by returning one list of id's instead of list of objects extra useless transformations, because no extra params except movie_id
func GetFavoriteMovieListDB(ctx context.Context, db *sql.DB, userID uuid.UUID) (models.FavoriteMovieList, error) {
	var getListSchema = `
		SELECT movie_id
		FROM favorite_movie
		WHERE user_id = $1
		`

	resRows, err := db.QueryContext(ctx, getListSchema, userID)
	if err != nil {
		return models.FavoriteMovieList{}, fmt.Errorf("get favorite_movie list for user: %w", err)
	}
	defer resRows.Close()

	favMovieList := models.FavoriteMovieList{UserID: userID}
	for resRows.Next() {
		select {
		case <-ctx.Done():
			return models.FavoriteMovieList{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			// Continue processing
		}

		favMovie := models.FavoriteMovieElem{}
		if err := resRows.Scan(&favMovie.MovieID); err != nil {
			return models.FavoriteMovieList{}, fmt.Errorf("reading favorite movie list: %w", err)
		}
		favMovieList.FavMovieList = append(favMovieList.FavMovieList, favMovie)
	}
	if err := resRows.Err(); err != nil {
		return models.FavoriteMovieList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return favMovieList, err
}
