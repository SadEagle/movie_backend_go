package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	"net/http"

	"movie_backend_go/db/sqlc"
	"movie_backend_go/internal/crudl"
	"movie_backend_go/internal/handlers/reqmodel"
	"movie_backend_go/pkg/auth"
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
		ho.Logger.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	favMovieList, err := ho.QuerierDB.GetUserFavoriteList(ctx, userID)
	if err != nil {
		ho.Logger.Printf("get user's favorite movie list from db: %v", err)
		http.Error(rw, "Can't get user's favorite movie list", http.StatusBadRequest)
		return
	}

	favUserListResp := reqmodel.UserFavoriteListResponse{UserID: userID, FavoriteMovieIDs: favMovieList}
	writeResponseBody(rw, favUserListResp, "user's favorite movie list")
}

// @Summary      Get my user favorite list
// @Description  Get current user favorite list
// @Tags         favorite, user
// @Accept       json
// @Produce      json
// @Security 		 OAuth2Password
// @Success      200  {object}  reqmodel.UserFavoriteListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/my/favorite [get]
func (ho *HandlerObj) GetMyUserFavoriteListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	favMovieList, err := ho.QuerierDB.GetUserFavoriteList(ctx, userTokenData.UserID)
	if err != nil {
		ho.Logger.Printf("get user's favorite movie list from db: %v", err)
		http.Error(rw, "Can't get user's favorite movie list", http.StatusBadRequest)
		return
	}

	favUserListResp := reqmodel.UserFavoriteListResponse{UserID: userTokenData.UserID, FavoriteMovieIDs: favMovieList}
	writeResponseBody(rw, favUserListResp, "user's favorite movie list")
}

// @Summary      Get movie favorite list
// @Description  Get list users who marked this movie as favorite
// @Tags         favorite, movie
// @Accept       json
// @Produce      json
// @Param        movie_id   	path	string  true  "Movie ID"
// @Success      200  {object}  reqmodel.MovieFavoriteListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id}/favorite [get]
func (ho *HandlerObj) GetMovieFavoriteListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	favUserList, err := ho.QuerierDB.GetUserFavoriteList(ctx, movieID)
	if err != nil {
		ho.Logger.Printf("get movie favorite list from db: %v", err)
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
// @Param        user_id		query		string	true	"User ID"
// @Param        movie_id   query		string	true	"Movie ID"
// @Success      200  {object}  sqlc.Favorite
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /favorite [get]
func (ho *HandlerObj) GetFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	movieIDStr := r.URL.Query().Get("movie_id")
	userIDStr := r.URL.Query().Get("user_id")
	if movieIDStr == "" || userIDStr == "" {
		ho.Logger.Println("Movie ID or User ID wasn't found")
		http.Error(rw, "Movie ID or User ID wasn't found", http.StatusBadRequest)
		return
	}
	var movieID, userID pgtype.UUID
	errMovie := movieID.Scan(movieID)
	errUser := userID.Scan(userID)
	if errMovie != nil || errUser != nil {
		ho.Logger.Println("Movie ID or User ID contain wrong style")
		http.Error(rw, "Movie ID or User ID contain wrong style", http.StatusBadRequest)
		return
	}

	favGet := sqlc.GetFavoriteParams{UserID: userID, MovieID: movieID}

	favorite, err := ho.QuerierDB.GetFavorite(ctx, favGet)
	if err != nil {
		ho.Logger.Printf("get favorite from db: %v", err)
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
// @Security 		 OAuth2Password
// @Param        request   body		reqmodel.FavoriteCreateRequest	true	"Comment update data"
// @Success      200  {object}  sqlc.Favorite
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /favorite [post]
func (ho *HandlerObj) CreateFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var favReq reqmodel.FavoriteCreateRequest
	err := decoder.Decode(&favReq)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	favCreate := sqlc.CreateFavoriteParams{UserID: userTokenData.UserID, MovieID: favReq.MovieID}

	favMovie, err := crudl.CreateFavorite(ctx, ho.QuerierDB, favCreate)

	if err != nil {
		ho.Logger.Printf("create movie favorite: %v", err)
		http.Error(rw, "Can't create favorite", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, favMovie, "favorite")
}

// @Summary		Delete favorite
// @Tags    	favorite, admin
// @Accept    json
// @Produce   json
// @Security  OAuth2Password
// @Param			request   body		reqmodel.FavoriteDeleteRequest	true	"Comment delete data"
// @Success   204
// @Failure   401  {object}  map[string]string
// @Failure   404  {object}  map[string]string
// @Failure   500  {object}  map[string]string
// @Router    /favorite [delete]
func (ho *HandlerObj) DeleteFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var favReq reqmodel.FavoriteDeleteRequest
	err := decoder.Decode(&favReq)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	if !userTokenData.IsAdmin {
		ho.Logger.Println(err)
		http.Error(rw, "Need to be admin for this action", http.StatusUnauthorized)
	}

	favDelete := sqlc.DeleteFavoriteParams{UserID: favReq.UserID, MovieID: favReq.MovieID}

	if err := crudl.DeleteFavorite(ctx, ho.QuerierDB, favDelete); err != nil {
		ho.Logger.Printf("Delete favorite: %v", err)
		http.Error(rw, "Can't delete favorite", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// @Summary		Delete my favorite
// @Tags    	favorite, user
// @Accept    json
// @Produce   json
// @Security  OAuth2Password
// @Param			request   body		reqmodel.MyFavoriteDeleteRequest	true	"Comment delete data"
// @Success   204
// @Failure   404  {object}  map[string]string
// @Failure   500  {object}  map[string]string
// @Router    /favorite/my [delete]
func (ho *HandlerObj) DeleteMyFavoriteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var favReq reqmodel.MyFavoriteDeleteRequest
	err := decoder.Decode(&favReq)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	favDelete := sqlc.DeleteFavoriteParams{UserID: userTokenData.UserID, MovieID: favReq.MovieID}

	if err := crudl.DeleteFavorite(ctx, ho.QuerierDB, favDelete); err != nil {
		ho.Logger.Printf("Delete favorite: %v", err)
		http.Error(rw, "Can't delete favorite", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
