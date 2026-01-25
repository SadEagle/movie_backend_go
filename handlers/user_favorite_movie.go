package handlers

import (
	"context"
	"fmt"
	"movie_backend_go/crudl"
	"movie_backend_go/reqmodel"
	"net/http"

	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary      Get favorite movie list
// @Description  Get user's favorite movie list
// @Tags         favorite_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Success      200  {object}  reqmodel.FavoriteMovieListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie [get]
func (ho *HandlerObj) GetFavoriteMovieListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	favMovieList, err := ho.DBPool.GetMovieFavoriteList(ctx, userID)
	if err != nil {
		ho.Log.Println(fmt.Errorf("get movie favorite list from db: %w", err))
		http.Error(rw, "Can't get user's favorite movie list", http.StatusBadRequest)
		return
	}

	favMovieListResp := reqmodel.FavoriteMovieListResponse{UserID: userID, FavoriteMovieIDs: favMovieList}
	writeResponseBody(rw, favMovieListResp, "favorite_move")
}

// @Summary      Add favorite movie
// @Description  Add movie to user's favorite
// @Tags         favorite_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      200  {object}  sqlc.FavoriteMovie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie/{movie_id} [post]
func (ho *HandlerObj) CreateMovieFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}
	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}
	favMovieCreate := sqlc.CreateMovieFavoriteParams{UserID: userID, MovieID: movieID}

	favMovie, err := crudl.CreateMovieFavorite(ctx, ho.DBPool, favMovieCreate)
	if err != nil {
		ho.Log.Println(fmt.Errorf("create movie favorite: %w", err))
		http.Error(rw, "Can't create favorite movie", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, favMovie, "favorite movie")
}

// @Summary      Delete favorite movie
// @Description  Delete favorite movie
// @Tags         favorite_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie/{movie_id} [delete]
func (ho *HandlerObj) DeleteFavoriteMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}
	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}
	favMovieDelete := sqlc.DeleteMovieFavoriteParams{UserID: userID, MovieID: movieID}

	if err := crudl.DeleteMovieFavorite(ctx, ho.DBPool, favMovieDelete); err != nil {
		ho.Log.Println(fmt.Errorf("Delete favorite movie: %w", err))
		http.Error(rw, "Can't delete favorite movie", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
