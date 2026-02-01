package handlers

import (
	"context"
	"encoding/json"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/db/sqlc"
	"movie_backend_go/reqmodel"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary 		Get movie comments list
// @Description	Get comment list for certain movie
// @Tags        comment, movie
// @Accept      json
// @Produce     json
// @Param       movie_id   path		string	true	"Movie ID"
// @Success     200		{object}	reqmodel.MovieCommentListResponse
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /movie/{movie_id}/comment [get]
func (ho *HandlerObj) GetMovieCommentListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var movieID pgtype.UUID
	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Printf("proceed path parameter: %v", err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	movieCommentList, err := crudl.GetMovieCommentList(ctx, ho.DBPool, movieID)
	if err != nil {
		ho.Logger.Printf("proceed getting movie comment list: %v", err)
		http.Error(rw, "Can't get movie comment list", http.StatusNotFound)
		return
	}
	movieCommentListResp := reqmodel.MovieCommentListResponse{MovieID: movieID, MovieCommentList: movieCommentList}
	writeResponseBody(rw, movieCommentListResp, "movie comment list")
}

// @Summary 		Get user comments list
// @Description	Get comment list for certain user
// @Tags        comment, user
// @Accept      json
// @Produce     json
// @Param       user_id   path		string	true	"User ID"
// @Success     200		{object}	reqmodel.UserCommentListResponse
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /user/{user_id}/comment [get]
func (ho *HandlerObj) GetUserCommentListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Logger.Printf("proceed path parameter: %v", err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	userCommentList, err := crudl.GetUserCommentList(ctx, ho.DBPool, userID)
	if err != nil {
		ho.Logger.Printf("proceed getting user comment list: %v", err)
		http.Error(rw, "Can't get user comment list", http.StatusNotFound)
		return
	}
	userCommentListResp := reqmodel.UserCommentListResponse{UserID: userID, UserCommentList: userCommentList}
	writeResponseBody(rw, userCommentListResp, "user comment list")
}

// @Summary 		Get comment
// @Description	Get comment by id
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       comment_id   path		string	true	"Comment ID"
// @Success     200		{object}	sqlc.Comment
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /comment/{comment_id} [get]
func (ho *HandlerObj) GetCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var commentID pgtype.UUID
	if err := commentID.Scan(r.PathValue("comment_id")); err != nil {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Requested comment id should contain uuid style", http.StatusBadRequest)
		return
	}

	comment, err := crudl.GetComment(ctx, ho.DBPool, commentID)
	if err != nil {
		ho.Logger.Printf("proceed getting comment comment list: %v", err)
		http.Error(rw, "Can't get user comment list", http.StatusNotFound)
		return
	}
	writeResponseBody(rw, comment, "user comment list")
}

// @Summary 		Create comments
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       request   body	reqmodel.CommentCreateRequest	true	"Comment create data"
// @Success     200		{object}	sqlc.Comment
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /comment [post]
func (ho *HandlerObj) CreateMovieCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var commentReq reqmodel.CommentCreateRequest
	err := decoder.Decode(&commentReq)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	commentCreate := sqlc.CreateCommentParams{Text: commentReq.Text}
	comment, err := crudl.CreateComment(ctx, ho.DBPool, commentCreate)
	if err != nil {
		ho.Logger.Printf("proceed creating comment: %v", err)
		http.Error(rw, "Can't create movie comment", http.StatusNotFound)
		return
	}
	writeResponseBody(rw, comment, "comment")
}

// @Summary 		Update comments
// @Tags        comment
// @Accept      json
// @Produce     json
// @Param       comment_id   path		string	true	"movie ID"
// @Param       request   body		reqmodel.CommentRequest	true	"Comment update data"
// @Success     200		{object}	sqlc.Comment
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /comment/{comment_id} [patch]
func (ho *HandlerObj) UpdateCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var commentID pgtype.UUID
	if err := commentID.Scan(r.PathValue("comment_id")); err != nil {
		ho.Logger.Printf("proceed path parameter: %v", err)
		http.Error(rw, "Requested comment id should contain uuid style", http.StatusBadRequest)
		return
	}

	var commentReq reqmodel.CommentRequest
	err := decoder.Decode(&commentReq)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	movieCommentUpdate := sqlc.UpdateCommentParams{ID: commentID, Text: commentReq.Text}
	movieComment, err := crudl.UpdateMovieComment(ctx, ho.DBPool, movieCommentUpdate)
	if err != nil {
		ho.Logger.Printf("proceed update comment: %v", err)
		http.Error(rw, "Can't update comment", http.StatusNotFound)
		return
	}
	writeResponseBody(rw, movieComment, "comment")
}

// @Summary      Delete comment
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        comment_id 	path	string 	true	"Movie ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comment/{comment_id} [delete]
func (ho *HandlerObj) DeleteMovieCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var commentID pgtype.UUID
	if err := commentID.Scan(r.PathValue("comment_id")); err != nil {
		ho.Logger.Printf("proceed path parameter: %v", err)
		http.Error(rw, "Requested comment id should contain uuid style", http.StatusBadRequest)
		return
	}

	if err := crudl.DeleteComment(ctx, ho.DBPool, commentID); err != nil {
		ho.Logger.Println("proceed delete movie comment request")
		http.Error(rw, "Can't delete movie comment", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
