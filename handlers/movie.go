package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"movie_backend_go/crudl"
	"movie_backend_go/models"

	"net/http"
)

// @Summary      Show movie
// @Description  Get movie by id
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Success      200  {object}  models.MovieResponse
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [get]
func GetMovieHandlerMake(db *sql.DB) http.HandlerFunc {
	GetMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		movie, err := crudl.GetMovieDB(db, id)
		if err != nil {
			log.Println(err)
			http.Error(rw, fmt.Sprintf("Can't get movie id: %s\n", id), 404)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		movieByte := movie.ToResponse()
		writeResponseBody(rw, movieByte, "movie")

	}
	return GetMovieHandler
}

// @Summary      Update movie
// @Description  Update movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Param        request 		body	models.UpdateMovieRequest  true  "Movie creation data"
// @Success      200  {object}  models.MovieResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [patch]
func UpdateMovieHandlerMake(db *sql.DB) http.HandlerFunc {
	UpdateMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Strict parsing

		var updateMovie models.UpdateMovieRequest
		err := decoder.Decode(&updateMovie)
		if err != nil && err != io.EOF {
			log.Println(err)
			http.Error(rw, "Can't proceed body request", 400)
			return
		}
		id := r.PathValue("id")

		movie, err := crudl.UpdateMovieDB(db, updateMovie, id)
		if err != nil {
			log.Println(err)
			http.Error(rw, "Can't update movie", 404)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		movieByte := movie.ToResponse()
		writeResponseBody(rw, movieByte, "movie")
	}
	return UpdateMovieHandler
}

// @Summary      Create movie
// @Description  Create movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        request 		body	models.CreateMovieRequest  true  "Movie creation data"
// @Success      201  {object}  models.MovieResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie [post]
func CreateMovieHandlerMake(db *sql.DB) http.HandlerFunc {
	CreateMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Strict parsing

		var createMovie models.CreateMovieRequest
		err := decoder.Decode(&createMovie)
		if err != nil && err != io.EOF {
			log.Println(err)
			http.Error(rw, "Can't proceed body request", 400)
			return
		}

		movie, err := crudl.CreateMovieDB(db, createMovie)
		if err != nil {
			log.Println(err)
			http.Error(rw, "Can't create movie", 404)
			return
		}

		rw.WriteHeader(201) // 201 - Create
		rw.Header().Set("Content-Type", "application/json")
		movieByte := movie.ToResponse()
		writeResponseBody(rw, movieByte, "movie")
	}
	return CreateMovieHandler
}

// @Description  Delete movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Success      204  {object}  models.MovieResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [delete]
func DeleteMovieHandlerMake(db *sql.DB) http.HandlerFunc {
	DeleteMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		err := crudl.DeleteMovieDB(db, id)
		if err != nil {
			log.Println(err)
			http.Error(rw, fmt.Sprintf("Can't delete movie id: %s", id), 404)
			return
		}
		rw.WriteHeader(204) // 204 - Success without returning body
	}
	return DeleteMovieHandler
}
