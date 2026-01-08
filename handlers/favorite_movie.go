package handlers

import (
	"context"
	"fmt"
	"movie_backend_go/crudl"
	"net/http"
	"time"
)

// @Description  Get user favorite_movie list
// @Tags         favorite_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Success      200  {object}  models.FavoriteMovieList
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie [get]
func (ho *HandlerObj) GetFavoriteMovieListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	user_id := r.PathValue("user_id")
	favMovieList, err := crudl.GetFavoriteMovieListDB(ctx, ho.DB, user_id)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't get favorite_movie list of user: %s", user_id), http.StatusNotFound)
		return
	}
	writeResponseBody(rw, favMovieList, "favorite_move")
}

// @Description  Add favorite movie
// @Tags         favorite_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      200  {object}  models.FavoriteMovie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie/{movie_id} [post]
func (ho *HandlerObj) AddFavoriteMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	userID := r.PathValue("user_id")
	movieID := r.PathValue("movie_id")
	favMovie, err := crudl.AddFavoriteMovieDB(ctx, ho.DB, userID, movieID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't add favorite movie user_id: %s, movie_id: %s", userID, userID), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, favMovie, "favorite movie")
}

// @Description  Delete favorite movie
// @Tags         favorite_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      204  {object}  models.FavoriteMovie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie/{movie_id} [delete]
func (ho *HandlerObj) DeleteFavoriteMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	user_id := r.PathValue("user_id")
	movie_id := r.PathValue("movie_id")
	err := crudl.DeleteFavoriteMovieDB(ctx, ho.DB, user_id, movie_id)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't delete favorite movie user_id: %s, movie_id: %s", user_id, user_id), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
