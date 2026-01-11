package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/models"
	"time"

	"net/http"

	"github.com/google/uuid"
)

// @Summary      Show movie
// @Description  Get movie by id
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Success      200  {object}  models.Movie
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [get]
func (ho *HandlerObj) GetMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	movieID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}
	movie, err := crudl.GetMovieByIDDB(ctx, ho.DB, movieID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't get movie id: %s\n", movieID), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	writeResponseBody(rw, movie, "movie")

}

// @Summary      Show movie list
// @Description  Get movie list
// @Tags         movie
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.MovieList
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie [get]
func (ho *HandlerObj) GetMovieListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	movieList, err := crudl.GetMovieListDB(ctx, ho.DB)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't get movie list", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	writeResponseBody(rw, movieList, "movie")

}

// @Summary      Update movie
// @Description  Update movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Param        request 		body	models.UpdateMovieRequest  true  "Movie creation data"
// @Success      200  {object}  models.Movie
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [patch]
func (ho *HandlerObj) UpdateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var updateMovie models.UpdateMovieRequest
	err := decoder.Decode(&updateMovie)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	movieID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	movie, err := crudl.UpdateMovieDB(ctx, ho.DB, updateMovie, movieID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't update movie", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	writeResponseBody(rw, movie, "movie")
}

// @Summary      Create movie
// @Description  Create movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        request 		body	models.CreateMovieRequest  true  "Movie creation data"
// @Success      201  {object}  models.Movie
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie [post]
func (ho *HandlerObj) CreateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict parsing

	var createMovie models.CreateMovieRequest
	err := decoder.Decode(&createMovie)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	movie, err := crudl.CreateMovieDB(ctx, ho.DB, createMovie)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't create movie", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	writeResponseBody(rw, movie, "movie list")
}

// @Description  Delete movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Success      204  {object}  models.Movie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [delete]
func (ho *HandlerObj) DeleteMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	movieID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	err = crudl.DeleteMovieDB(ctx, ho.DB, movieID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't delete movie id: %s", movieID), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
