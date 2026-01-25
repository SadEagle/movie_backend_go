package handlers

import (
	"context"
	"fmt"
	"movie_backend_go/crudl"
	"movie_backend_go/reqmodel"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary      Get user comment list
// @Description  Get user's comment list
// @Tags         user, movie_comment
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Success      200  {object}  reqmodel.UserCommentListResp
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/comments [get]
func (ho *HandlerObj) GetUserCommentListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	userCommentList, err := crudl.GetUserCommentList(ctx, ho.DBPool, userID)
	if err != nil {
		ho.Log.Println(fmt.Errorf("getting user comment list: %w", err))
		http.Error(rw, "Can't get user's comment list", http.StatusBadRequest)
		return
	}

	movieCommentListResp := reqmodel.UserCommentListResp{UserID: userID, UserCommentList: userCommentList}
	writeResponseBody(rw, movieCommentListResp, "user comment list")
}
