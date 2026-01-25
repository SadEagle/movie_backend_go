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

// @Summary 		Get movie comments list
// @Description	Get movie comments list
// @Tags        movie, movie_comment
// @Accept      json
// @Produce     json
// @Param       movie_id   path		string	true	"movie ID"
// @Success     200		{object}	reqmodel.MovieCommentListResp
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /movie/{movie_id}/comment [get]
func (ho *HandlerObj) GetMovieCommentListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(fmt.Errorf("proceed body request: %w", err))
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	movieCommentList, err := crudl.GetMovieCommentList(ctx, ho.DBPool, movieID)
	if err != nil {
		ho.Log.Println(fmt.Errorf("proceed getting movie comment list: %w", err))
		http.Error(rw, "Can't get movie comment list", http.StatusNotFound)
		return
	}
	movieCommentListResp := reqmodel.MovieCommentListResp{MovieID: movieID, MovieCommentList: movieCommentList}
	writeResponseBody(rw, movieCommentListResp, "movie comment list")
}

// @Summary 		Create movie comments list
// @Description	Create movie comments list
// @Tags        movie, movie_comment
// @Accept      json
// @Produce     json
// @Param       movie_id   path		string	true	"movie ID"
// @Success     200		{object}	sqlc.MovieComment
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /movie/{movie_id}/comment [post]
func (ho *HandlerObj) CreateMovieCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(fmt.Errorf("proceed body request: %w", err))
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	var movieCommentReq reqmodel.MovieCommentRequest
	err := decoder.Decode(&movieCommentReq)
	if err != nil && err != io.EOF {
		ho.Log.Println(fmt.Errorf("proceed body request: %w", err))
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	movieCommentCreate := sqlc.CreateMovieCommentParams{UserID: movieCommentReq.UserID, MovieID: movieID, Text: movieCommentReq.Text}

	movieComment, err := crudl.CreateMovieComment(ctx, ho.DBPool, movieCommentCreate)
	if err != nil {
		ho.Log.Println(fmt.Errorf("proceed creating movie comment: %w", err))
		http.Error(rw, "Can't create movie comment", http.StatusNotFound)
		return
	}
	writeResponseBody(rw, movieComment, "movie comment")
}

// @Summary 		Update movie comments list
// @Description	Update movie comments list
// @Tags        movie, movie_comment
// @Accept      json
// @Produce     json
// @Param       movie_id   path		string	true	"movie ID"
// @Success     200		{object}	sqlc.MovieComment
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /movie/{movie_id}/comment [patch]
func (ho *HandlerObj) UpdateMovieCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(fmt.Errorf("proceed body request: %w", err))
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	var movieCommentReq reqmodel.MovieCommentRequest
	err := decoder.Decode(&movieCommentReq)
	if err != nil && err != io.EOF {
		ho.Log.Println(fmt.Errorf("proceed body request: %w", err))
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	movieCommentUpdate := sqlc.UpdateMovieCommentParams{UserID: movieCommentReq.UserID, MovieID: movieID, Text: movieCommentReq.Text}

	movieComment, err := crudl.UpdateMovieComment(ctx, ho.DBPool, movieCommentUpdate)
	if err != nil {
		ho.Log.Println(fmt.Errorf("proceed update movie comment: %w", err))
		http.Error(rw, "Can't update movie comment", http.StatusNotFound)
		return
	}
	writeResponseBody(rw, movieComment, "movie comment")
}

// @Summary      Delete movie comment
// @Description  Delete movie comment
// @Tags         movie, movie_comment
// @Accept       json
// @Produce      json
// @Param        movie_id   path      string  true  "Movie ID"
// @Param        user_id   	path      string  true  "User ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{movie_id} [delete]
func (ho *HandlerObj) DeleteMovieCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}
	var userID pgtype.UUID
	if err := movieID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	movieCommentDelete := sqlc.DeleteMovieCommentParams{UserID: userID, MovieID: movieID}

	if err := crudl.DeleteMovieComment(ctx, ho.DBPool, movieCommentDelete); err != nil {
		ho.Log.Println(fmt.Errorf("proceed delete movie comment request"))
		http.Error(rw, "Can't delete movie comment", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
