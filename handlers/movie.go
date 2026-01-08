package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"movie_backend_go/crudl"
	"movie_backend_go/models"
	"time"

	"net/http"
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
	movieID := r.PathValue("id")
	movie, err := crudl.GetMovieDB(ctx, ho.DB, movieID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't get movie id: %s\n", movieID), 404)
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
		http.Error(rw, "Can't get movie list", 500)
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
	decoder.DisallowUnknownFields() // Strict parsing

	var updateMovie models.UpdateMovieRequest
	err := decoder.Decode(&updateMovie)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", 400)
		return
	}
	id := r.PathValue("id")

	movie, err := crudl.UpdateMovieDB(ctx, ho.DB, updateMovie, id)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't update movie", 404)
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
		http.Error(rw, "Can't proceed body request", 400)
		return
	}

	movie, err := crudl.CreateMovieDB(ctx, ho.DB, createMovie)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't create movie", 404)
		return
	}

	rw.WriteHeader(201) // 201 - Create
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
	id := r.PathValue("id")
	err := crudl.DeleteMovieDB(ctx, ho.DB, id)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't delete movie id: %s", id), 404)
		return
	}
	rw.WriteHeader(204) // 204 - Success without returning body
}
