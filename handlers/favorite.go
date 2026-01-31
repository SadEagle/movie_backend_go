package handlers

import (
	"context"
	"encoding/json"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/reqmodel"
	"net/http"

	"movie_backend_go/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary      Get user favorite list
// @Description  Get user's favorite movie list
// @Tags         favorite, user
// @Accept       json
// @Produce      json
// @Param        user_id   	path	string  true  "User ID"
// @Success      200  {object}  reqmodel.UserFavoriteListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite [get]
func (ho *HandlerObj) GetUserFavoriteListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	favMovieList, err := ho.DBPool.GetUserFavoriteList(ctx, userID)
	if err != nil {
		ho.Log.Printf("get user's favorite movie list from db: %v", err)
		http.Error(rw, "Can't get user's favorite movie list", http.StatusBadRequest)
		return
	}

	favUserListResp := reqmodel.UserFavoriteListResponse{UserID: userID, FavoriteMovieIDs: favMovieList}
	writeResponseBody(rw, favUserListResp, "user's favorite movie list")

}

// @Summary      Get movie favorite list
// @Description  Get list users who marked this movie as favorite
// @Tags         favorite, movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path	string  true  "User ID"
// @Success      200  {object}  reqmodel.MovieFavoriteListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{user_id}/favorite [get]
func (ho *HandlerObj) GetMovieFavoriteListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	favUserList, err := ho.DBPool.GetUserFavoriteList(ctx, movieID)
	if err != nil {
		ho.Log.Printf("get movie favorite list from db: %v", err)
		http.Error(rw, "Can't get user's favorite movie list", http.StatusBadRequest)
		return
	}

	favMovieListResp := reqmodel.MovieFavoriteListResponse{MovieID: movieID, FavoriteUserIDs: favUserList}
	writeResponseBody(rw, favMovieListResp, "movie's favorite user list")
}

// @Summary      Get favorite
// @Tags         favorite
// @Accept       json
// @Produce      json
// @Param       request   body		reqmodel.FavoriteRequest	true	"Comment update data"
// @Success      200  {object}  sqlc.Favorite
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /favorite [get]
func (ho *HandlerObj) GetFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var favReq reqmodel.FavoriteRequest
	err := decoder.Decode(&favReq)
	if err != nil && err != io.EOF {
		ho.Log.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	favGet := sqlc.GetFavoriteParams{UserID: favReq.UserID, MovieID: favReq.MovieID}

	favorite, err := ho.DBPool.GetFavorite(ctx, favGet)
	if err != nil {
		ho.Log.Printf("get favorite from db: %v", err)
		http.Error(rw, "Can't get favorite", http.StatusBadRequest)
		return
	}

	writeResponseBody(rw, favorite, "favorite")
}

// @Summary      Add favorite movie
// @Description  Add movie to user's favorite
// @Tags         favorite
// @Accept       json
// @Produce      json
// @Param       request   body		reqmodel.FavoriteRequest	true	"Comment update data"
// @Success      200  {object}  sqlc.Favorite
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /favorite [post]
func (ho *HandlerObj) CreateFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var favReq reqmodel.FavoriteRequest
	err := decoder.Decode(&favReq)
	if err != nil && err != io.EOF {
		ho.Log.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	favCreate := sqlc.CreateFavoriteParams{UserID: favReq.UserID, MovieID: favReq.MovieID}

	favMovie, err := crudl.CreateFavorite(ctx, ho.DBPool, favCreate)
	if err != nil {
		ho.Log.Printf("create movie favorite: %v", err)
		http.Error(rw, "Can't create favorite", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, favMovie, "favorite")
}

// @Summary		Delete favorite
// @Tags    	favorite
// @Accept    json
// @Produce   json
// @Param			request   body		reqmodel.FavoriteRequest	true	"Comment delete data"
// @Success   204
// @Failure   404  {object}  map[string]string
// @Failure   500  {object}  map[string]string
// @Router    /favorite [delete]
func (ho *HandlerObj) DeleteFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var favReq reqmodel.FavoriteRequest
	err := decoder.Decode(&favReq)
	if err != nil && err != io.EOF {
		ho.Log.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	favDelete := sqlc.DeleteFavoriteParams{UserID: favReq.UserID, MovieID: favReq.MovieID}

	if err := crudl.DeleteFavorite(ctx, ho.DBPool, favDelete); err != nil {
		ho.Log.Printf("Delete favorite: %v", err)
		http.Error(rw, "Can't delete favorite", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
