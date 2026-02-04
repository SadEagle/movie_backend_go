package handlers

import (
	"context"
	"encoding/json"
	"io"
	"movie_backend_go/db/sqlc"
	"movie_backend_go/internal/auth"
	"movie_backend_go/internal/crudl"
	"movie_backend_go/internal/handlers/reqmodel"
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

	movieCommentList, err := crudl.GetMovieCommentList(ctx, ho.QuerierDB, movieID)
	if err != nil {
		ho.Logger.Printf("proceed getting movie comment list: %v", err)
		http.Error(rw, "Can't get movie comment list", http.StatusNotFound)
		return
	}
	movieCommentListResp := reqmodel.MovieCommentListResponse{MovieID: movieID, MovieCommentList: movieCommentList}
	writeResponseBody(rw, movieCommentListResp, "movie comment list")
}

// @Summary 		Get my user comments list
// @Description	Get current user comment list
// @Tags        comment, user
// @Accept      json
// @Produce     json
// @Security 		OAuth2Password
// @Success     200		{object}	reqmodel.UserCommentListResponse
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /user/my/comment [get]
func (ho *HandlerObj) GetMyUserCommentListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	userCommentList, err := crudl.GetUserCommentList(ctx, ho.QuerierDB, userTokenData.UserID)
	if err != nil {
		ho.Logger.Printf("proceed getting user comment list: %v", err)
		http.Error(rw, "Can't get user comment list", http.StatusNotFound)
		return
	}
	userCommentListResp := reqmodel.UserCommentListResponse{UserID: userTokenData.UserID, UserCommentList: userCommentList}
	writeResponseBody(rw, userCommentListResp, "user comment list")
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

	userCommentList, err := crudl.GetUserCommentList(ctx, ho.QuerierDB, userID)
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

	comment, err := crudl.GetComment(ctx, ho.QuerierDB, commentID)
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
// @Security 		OAuth2Password
// @Param       request   body	reqmodel.CommentCreateRequest	true	"Comment create data"
// @Success     200		{object}	sqlc.Comment
// @Failure     401  	{object}  map[string]string
// @Failure     404  	{object}  map[string]string
// @Failure     500  	{object}  map[string]string
// @Router      /comment [post]
func (ho *HandlerObj) CreateCommentHandler(rw http.ResponseWriter, r *http.Request) {
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

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	// Verify
	commentCreate := sqlc.CreateCommentParams{UserID: userTokenData.UserID, MovieID: commentReq.MovieID, Text: commentReq.Text}

	comment, err := crudl.CreateComment(ctx, ho.QuerierDB, commentCreate)
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
// @Security 		OAuth2Password
// @Param       comment_id   path		string	true	"movie ID"
// @Param       request   body		reqmodel.CommentUpdateRequest	true	"Comment update data"
// @Success     200		{object}	sqlc.Comment
// @Failure     401  	{object}  map[string]string
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

	var commentReq reqmodel.CommentUpdateRequest
	err := decoder.Decode(&commentReq)
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

	// Verify
	commentData, err := crudl.GetComment(ctx, ho.QuerierDB, commentID)
	if err != nil {
		ho.Logger.Fatalf("searching comment by id - %s: %v", commentID, err)
		http.Error(rw, "Can't find comment with current id", http.StatusBadRequest)
		return
	}
	if commentData.UserID != userTokenData.UserID {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	commentUpdate := sqlc.UpdateCommentParams{ID: commentID, Text: commentReq.Text}
	comment, err := crudl.UpdateMovieComment(ctx, ho.QuerierDB, commentUpdate)
	if err != nil {
		ho.Logger.Printf("proceed update comment: %v", err)
		http.Error(rw, "Can't update comment", http.StatusNotFound)
		return
	}
	writeResponseBody(rw, comment, "comment")
}

// @Summary      Delete comment
// @Tags         comment, admin, user
// @Accept       json
// @Produce      json
// @Security 		 OAuth2Password
// @Param        comment_id 	path	string 	true	"Movie ID"
// @Success      204
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /comment/{comment_id} [delete]
func (ho *HandlerObj) DeleteCommentHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var commentID pgtype.UUID
	if err := commentID.Scan(r.PathValue("comment_id")); err != nil {
		ho.Logger.Printf("proceed path parameter: %v", err)
		http.Error(rw, "Requested comment id should contain uuid style", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	// Verify token user or admin is owner
	if !userTokenData.IsAdmin {
		commentData, err := crudl.GetComment(ctx, ho.QuerierDB, commentID)
		if err != nil {
			// TODO: add extra information like commentID in error Logger like this
			ho.Logger.Fatalf("searching comment: %v", err)
			http.Error(rw, "Can't find comment with current id", http.StatusBadRequest)
			return
		}
		if commentData.UserID != userTokenData.UserID {
			ho.Logger.Println("Unauthorized user")
			http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
			return
		}
	}

	if err := crudl.DeleteComment(ctx, ho.QuerierDB, commentID); err != nil {
		ho.Logger.Println("proceed delete movie comment request")
		http.Error(rw, "Can't delete movie comment", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
