package handlers

import (
	"context"
	"encoding/json"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/db/sqlc"

	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"movie_backend_go/reqmodel"
)

// @Summary      Get movie list
// @Description  Get all movie list
// @Tags         movie
// @Accept       json
// @Produce      json
// @Success      200  {object}  reqmodel.MovieListResponse
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie [get]
func (ho *HandlerObj) GetMovieListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	movieList, err := crudl.GetMovieList(ctx, ho.DBPool)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Can't get movie list", http.StatusBadRequest)
		return
	}

	movieListResponse := reqmodel.MovieListResponse{MovieList: movieList}
	writeResponseBody(rw, movieListResponse, "movie")
}

// @Summary      Get movie
// @Description  Get movie by id
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      200  {object}  sqlc.GetMovieByIDRow
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id} [get]
func (ho *HandlerObj) GetMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}
	movie, err := crudl.GetMovieByID(ctx, ho.DBPool, movieID)
	if err != nil {
		ho.Logger.Printf("get movie by id %v: %v", movieID, err)
		http.Error(rw, "Can't get movie by id", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, movie, "movie")
}

// @Summary      Update movie
// @Description  Update movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        movie_id   path      string  true  "Movie ID"
// @Param        request 		body	reqmodel.MovieRequest  true  "Movie creation data"
// @Success      200  {object}  sqlc.Movie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id} [patch]
func (ho *HandlerObj) UpdateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var updateMovie reqmodel.MovieUpdateRequest
	err := decoder.Decode(&updateMovie)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	movieUpdate := sqlc.UpdateMovieParams{ID: movieID, Title: updateMovie.Title}
	movie, err := ho.DBPool.UpdateMovie(ctx, movieUpdate)
	if err != nil {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, movie, "movie")
}

// @Summary      Create movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        request 		body	reqmodel.MovieRequest  true  "Movie creation data"
// @Success      201  {object}  sqlc.Movie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie [post]
func (ho *HandlerObj) CreateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var movieCreate reqmodel.MovieRequest
	err := decoder.Decode(&movieCreate)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	movie, err := crudl.CreateMovie(ctx, ho.DBPool, movieCreate.Title)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Can't create movie", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, movie, "movie")
}

// @Summary      Delete movie
// @Tags         movie
// @Accept       json
// @Produce      json
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id} [delete]
func (ho *HandlerObj) DeleteMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	if err := crudl.DeleteMovie(ctx, ho.DBPool, movieID); err != nil {
		ho.Logger.Println("proceed delete movie request")
		http.Error(rw, "Can't delete movie", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
