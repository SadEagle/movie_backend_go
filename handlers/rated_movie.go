package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// @Description  Get user rated_movie list
// @Tags         rated_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Success      200  {object}  models.RatedMovieList
// @Failure      404  {object}  map[string]string
// @Router       /user/{user_id}/rated_movie [get]
func (ho *HandlerObj) GetRatedMovieListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()

	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid", http.StatusBadRequest)
		return
	}

	ratedMovieList, err := crudl.GetRatedMovieListDB(ctx, ho.DB, userID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't get rated movie list of user: %s", userID), http.StatusNotFound)
		return
	}
	writeResponseBody(rw, ratedMovieList, "rated move list")
}

// @Description  Get user rated_movie list
// @Tags         rated_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        request   	body      models.RatedMovieListElem  true  "Rate movie data"
// @Success      200  {object}  models.RatedMovie
// @Failure      404  {object}  map[string]string
// @Router       /user/{user_id}/rated_movie [get]
func (ho *HandlerObj) AddRatedMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid", http.StatusBadRequest)
		return
	}

	var ratedMovieElem models.RatedMovieElem
	err = decoder.Decode(&ratedMovieElem)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	ratedMovieList, err := crudl.GetRatedMovieListDB(ctx, ho.DB, userID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't get rated movie list of user: %s", userID), http.StatusNotFound)
		return
	}
	writeResponseBody(rw, ratedMovieList, "rated move list")
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
func (ho *HandlerObj) DeleteRatedMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()

	userID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid", http.StatusBadRequest)
		return
	}
	movieID, err := uuid.Parse(r.PathValue("movie_id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid", http.StatusBadRequest)
		return
	}

	err = crudl.DeleteRatedMovieDB(ctx, ho.DB, userID, movieID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't delete rated movie user_id: %s, movie_id: %s", userID, userID), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// TODO: Add update
