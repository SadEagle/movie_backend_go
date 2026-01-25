package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/db/sqlc"
	"movie_backend_go/reqmodel"

	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary			 Get  movie rating
// @Description  Get movie rating by user id
// @Tags         rated_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Success      200  {object}  reqmodel.RatedMovieListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/rated_movie [get]
func (ho *HandlerObj) GetRatedMovieListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	ratedMovieList, err := crudl.GetMovieRatingList(ctx, ho.DBPool, userID)
	if err != nil {
		ho.Log.Println(fmt.Errorf("proceed rated movie list: %w", err))
		http.Error(rw, "Can't proceed rated movie list", http.StatusNotFound)
		return
	}
	ratedMovieListResponse := reqmodel.RatedMovieListResponse{UserID: userID, RatedMovieList: ratedMovieList}

	writeResponseBody(rw, ratedMovieListResponse, "rated move list")
}

// @Summary			 Rate movie
// @Description  Rate movie
// @Tags         rated_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        request   	body      reqmodel.RatedMovieRequest  true  "Rate movie data"
// @Success      200  {object}  sqlc.RatedMovie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/rated_movie [post]
func (ho *HandlerObj) CreateRatedMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	var ratedMovieRequest reqmodel.RatedMovieRequest
	err := decoder.Decode(&ratedMovieRequest)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	ratedMovieCreate := sqlc.CreateMovieRatingParams{UserID: userID, MovieID: ratedMovieRequest.MovieID, Rating: ratedMovieRequest.Rating}

	ratedMovie, err := crudl.CreateMovieRating(ctx, ho.DBPool, ratedMovieCreate)
	if err != nil {
		ho.Log.Println(fmt.Errorf("proceed rated movie creation: %w", err))
		http.Error(rw, "Can't create movie rating", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, ratedMovie, "rated move list")
}

// @Summary			 Update movie rating
// @Description  Update movie rating
// @Tags         rated_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        request   	body      reqmodel.RatedMovieRequest  true  "Updated rating"
// @Success      200  {object}  sqlc.RatedMovie
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/rated_movie [patch]
func (ho *HandlerObj) UpdateRatedMovieHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	var ratedMovieUpdateRequest reqmodel.RatedMovieUpdateRequest
	err := decoder.Decode(&ratedMovieUpdateRequest)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	ratedMovieUpdate := sqlc.UpdateMoveRatingParams{UserID: userID, MovieID: ratedMovieUpdateRequest.MovieID, Rating: ratedMovieUpdateRequest.Rating}

	ratedMovie, err := crudl.UpdateMovieRating(ctx, ho.DBPool, ratedMovieUpdate)
	if err != nil {
		ho.Log.Println(fmt.Errorf("proceed rated movie creation: %w", err))
		http.Error(rw, "Can't update movie rating", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, ratedMovie, "rated move list")
}

// @Summary			 Delete movie rating
// @Description  Delete certain movie rating
// @Tags         rated_movie, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/rated_movie/{movie_id} [delete]
func (ho *HandlerObj) DeleteRatedMovieHandler(rw http.ResponseWriter, r *http.Request) {
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

	movieRatingDelete := sqlc.DeleteMovieRatingParams{UserID: userID, MovieID: movieID}
	if err := crudl.DeleteMovieRating(ctx, ho.DBPool, movieRatingDelete); err != nil {
		ho.Log.Println(fmt.Errorf("proceed delete movie rating request"))
		http.Error(rw, "Can't delete movie rating", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
