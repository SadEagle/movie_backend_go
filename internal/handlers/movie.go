package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"movie_backend_go/db/sqlc"
	"movie_backend_go/internal/crudl"
	"movie_backend_go/internal/handlers/reqmodel"
	"movie_backend_go/pkg/auth"

	"github.com/jackc/pgx/v5/pgtype"
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

	movieList, err := crudl.GetMovieList(ctx, ho.QuerierDB)
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
// @Success      200  {object}  sqlc.Movie
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

	movie, err := crudl.GetMovie(ctx, ho.QuerierDB, movieID)
	if err != nil {
		ho.Logger.Printf("get movie by id %v: %v", movieID, err)
		http.Error(rw, "Can't get movie by id", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, movie, "movie")
}

// TODO: add GetMovieByTitle

// @Summary      Update movie
// @Description  Update movie
// @Tags         movie, admin
// @Accept       json
// @Produce      json
// @Security	 OAuth2Password
// @Param        movie_id   path      string  true  "Movie ID"
// @Param        request 		body	reqmodel.MovieUpdateRequest  true  "Movie creation data"
// @Success      200  {object}  sqlc.Movie
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id} [patch]
func (ho *HandlerObj) UpdateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var movieUpdateReq reqmodel.MovieUpdateRequest
	err := decoder.Decode(&movieUpdateReq)
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

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	// Verify
	if !userTokenData.IsAdmin {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	movieUpdate := sqlc.UpdateMovieParams{ID: movieID, Title: movieUpdateReq.Title}
	movie, err := ho.QuerierDB.UpdateMovie(ctx, movieUpdate)
	if err != nil {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, movie, "movie")
}

// @Summary      Create movie
// @Tags         movie, admin
// @Accept       json
// @Produce      json
// @Security	 OAuth2Password
// @Param        request 		body	reqmodel.MovieCreateRequest  true  "Movie creation data"
// @Success      201  {object}  sqlc.Movie
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie [post]
func (ho *HandlerObj) CreateMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var movieCreateReq reqmodel.MovieCreateRequest
	err := decoder.Decode(&movieCreateReq)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
		return
	}
	// Verify
	if !userTokenData.IsAdmin {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	movie, err := crudl.CreateMovie(ctx, ho.QuerierDB, movieCreateReq.Title)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Can't create movie", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, movie, "movie")
}

// @Summary      Delete movie
// @Tags         movie, admin
// @Accept       json
// @Produce      json
// @Security	 OAuth2Password
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      204
// @Failure      401 	{object}  map[string]string
// @Failure      404  	{object}  map[string]string
// @Failure      500  	{object}  map[string]string
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

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}
	// Verify
	if !userTokenData.IsAdmin {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	if err := crudl.DeleteMovie(ctx, ho.QuerierDB, movieID); err != nil {
		ho.Logger.Println("proceed delete movie request")
		http.Error(rw, "Can't delete movie", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
