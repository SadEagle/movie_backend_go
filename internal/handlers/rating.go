package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	"net/http"

	"movie_backend_go/db/sqlc"
	"movie_backend_go/internal/auth"
	"movie_backend_go/internal/crudl"
	"movie_backend_go/internal/handlers/reqmodel"
)

// @Summary			 Get user rating list
// @Description  Get user's rated movies
// @Tags         rating, user
// @Accept       json
// @Produce      json
// @Param        user_id 	path	string  true  "User ID"
// @Success      200  {object}  reqmodel.UserRatingListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/rating [get]
func (ho *HandlerObj) GetUserRatingListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	userRatingList, err := crudl.GetUserRatingList(ctx, ho.QuerierDB, userID)
	if err != nil {
		ho.Logger.Printf("proceed rated movie list: %v", err)
		http.Error(rw, "Can't proceed rated movie list", http.StatusNotFound)
		return
	}
	ratedMovieListResponse := reqmodel.UserRatingListResponse{UserID: userID, UserRatingList: userRatingList}

	writeResponseBody(rw, ratedMovieListResponse, "user rating list")
}

// @Summary			 Get movie rating list
// @Description  Get users who rated movie
// @Tags         rating, movie
// @Accept       json
// @Produce      json
// @Param        movie_id 	path	string  true  "Movie ID"
// @Success      200  {object}  reqmodel.MovieRatingListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id}/rating [get]
func (ho *HandlerObj) GetMovieRatingListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	movieRatingList, err := crudl.GetMovieRatingList(ctx, ho.QuerierDB, movieID)
	if err != nil {
		ho.Logger.Printf("proceed users who rate current movie: %v", err)
		http.Error(rw, "Can't proceed users who rate current movie", http.StatusNotFound)
		return
	}
	ratedMovieListResponse := reqmodel.MovieRatingListResponse{MovieID: movieID, MovieRatingList: movieRatingList}

	writeResponseBody(rw, ratedMovieListResponse, "movie rating list")
}

// @Summary			 Get rating
// @Tags         rating
// @Accept       json
// @Produce      json
// @Param        request   	body      reqmodel.RatingIDRequest  true  "Rating ID"
// @Success      200  {object}  sqlc.Rating
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
func (ho *HandlerObj) GetRatingHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var ratingReq reqmodel.RatingIDRequest
	err := decoder.Decode(&ratingReq)
	if err != nil && err != io.EOF {
		ho.Logger.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	ratingGet := sqlc.GetRatingParams{UserID: ratingReq.UserID, MovieID: ratingReq.MovieID}

	rating, err := crudl.GetRating(ctx, ho.QuerierDB, ratingGet)
	if err != nil {
		ho.Logger.Printf("proceed get rating: %v", err)
		http.Error(rw, "Can't proceed get rating", http.StatusNotFound)
		return
	}

	writeResponseBody(rw, rating, "rating")
}

// @Summary			 Create rating
// @Description  Rate movie by user
// @Tags         rating
// @Accept       json
// @Produce      json
// @Security	 	 OAuth2Password
// @Param        request   	body      reqmodel.RatingCreateRequest  true  "Rate movie data"
// @Success      200  {object}  sqlc.Rating
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /rating [post]
func (ho *HandlerObj) CreateRatingHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var ratingReq reqmodel.RatingCreateRequest
	err := decoder.Decode(&ratingReq)
	if err != nil && err != io.EOF {
		ho.Logger.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	ratingCreate := sqlc.CreateRatingParams{UserID: userTokenData.UserID, MovieID: ratingReq.MovieID, Rating: ratingReq.Rating}

	rating, err := crudl.CreateMovieRating(ctx, ho.QuerierDB, ratingCreate)
	if err != nil {
		ho.Logger.Printf("proceed rating creation: %v", err)
		http.Error(rw, "Can't create rating", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, rating, "rating")
}

// @Summary			 Update rating
// @Description  Update rating
// @Tags         rating
// @Accept       json
// @Produce      json
// @Security	 	 OAuth2Password
// @Param        request   	body      reqmodel.RatingUpdateRequest  true  "Updated rating"
// @Success      200  {object}  sqlc.Rating
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /rating [patch]
func (ho *HandlerObj) UpdateRatingHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var ratingUpdateReq reqmodel.RatingUpdateRequest
	err := decoder.Decode(&ratingUpdateReq)
	if err != nil && err != io.EOF {
		ho.Logger.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	ratingUpdate := sqlc.UpdateRatingParams{UserID: userTokenData.UserID, MovieID: ratingUpdateReq.MovieID, Rating: ratingUpdateReq.Rating}

	ratedMovie, err := crudl.UpdateMovieRating(ctx, ho.QuerierDB, ratingUpdate)
	if err != nil {
		ho.Logger.Printf("proceed rating update: %v", err)
		http.Error(rw, "Can't update rating", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, ratedMovie, "rating")
}

// @Summary			 Delete rating
// @Description  Delete certain rating
// @Tags         rating, admin
// @Accept       json
// @Produce      json
// @Security	 	 OAuth2Password
// @Success      204
// @Failure      401 {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /rating [delete]
func (ho *HandlerObj) DeleteRatingHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var ratingDeleteReq reqmodel.RatingIDRequest
	err := decoder.Decode(&ratingDeleteReq)
	if err != nil && err != io.EOF {
		ho.Logger.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}
	// Verify token user or admin is owner
	if !userTokenData.IsAdmin && userTokenData.UserID != ratingDeleteReq.UserID {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	ratingDelete := sqlc.DeleteRatingParams{UserID: ratingDeleteReq.UserID, MovieID: ratingDeleteReq.MovieID}
	if err := crudl.DeleteRating(ctx, ho.QuerierDB, ratingDelete); err != nil {
		ho.Logger.Printf("proceed delete rating request: %v", err)
		http.Error(rw, "Can't delete rating", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
