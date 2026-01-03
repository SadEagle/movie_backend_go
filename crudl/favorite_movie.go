package crudl

import (
	"database/sql"
	"fmt"
	"movie_backend_go/models"
)

func AddFavoriteMovieDB(db *sql.DB, user_id string, movie_id string) (models.FavoriteMovie, error) {

	var createShema = `
		INSERT INTO favorite_movie(user_id, movie_id)
		VALUES ($1, $2)
		RETURNING user_id, movie_id
		`
	res_create := db.QueryRow(createShema, user_id, movie_id)

	favoriteMovie := models.FavoriteMovie{}
	err := res_create.Scan(&favoriteMovie.UserID, &favoriteMovie.MovieID)
	if err != nil {
		return models.FavoriteMovie{}, fmt.Errorf("scanning favorite_movie adding: %w", err)
	}
	return favoriteMovie, nil
}

func DeleteFavoriteMovieDB(db *sql.DB, user_id string, movie_id string) error {
	var createShema = `
		DELETE FROM favorite_movie
		WHERE user_id = $1 AND movie_id = $2
		`
	res, err := db.Exec(createShema, user_id, movie_id)
	if err != nil {
		return fmt.Errorf("delete favorite_movie: %w", err)
	}
	return checkNonEmptyDeletion(res)
}

// Simplified version. Better option will be create non-response datatype and later convert to response one
func GetFavoriteMovieListDB(db *sql.DB, user_id string) (models.FavoriteMovieList, error) {
	var getListSchema = `
		SELECT user_id, movie_id
		FROM favorite_movie
		WHERE user_id = $1
		`
	rows, err := db.Query(getListSchema, user_id)
	if err != nil {
		return models.FavoriteMovieList{}, fmt.Errorf("get favorite_movie list for user: %w", err)
	}
	defer rows.Close()

	favMovieList := models.FavoriteMovieList{}
	for rows.Next() {
		var favMovie models.FavoriteMovie
		if err := rows.Scan(&favMovie.UserID, &favMovie.MovieID); err != nil {
			return models.FavoriteMovieList{}, fmt.Errorf("reading favorite movie list: %w", err)
		}
		favMovieList.FavMovieList = append(favMovieList.FavMovieList, favMovie)
	}
	if err := rows.Err(); err != nil {
		return models.FavoriteMovieList{}, fmt.Errorf("check for errors from iteration over rows: %w", err)
	}
	return favMovieList, err
}
