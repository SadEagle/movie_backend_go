package handlers

import (
	"context"
	"encoding/json"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/db/sqlc"
	"movie_backend_go/reqmodel"
	"movie_backend_go/secret"

	"github.com/jackc/pgx/v5/pgtype"

	"net/http"
)

// @Summary      Show user
// @Description  Get user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id   path      string  true  "User ID"
// @Success      200  {object}  sqlc.UserDatum
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id} [get]
func (ho *HandlerObj) GetUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}
	user, err := crudl.GetUserByID(ctx, ho.DBPool, userID)
	if err != nil {
		ho.Logger.Printf("proceed getting user: %v", err)
		http.Error(rw, "Can't proceed getting user", http.StatusBadRequest)
		return
	}

	writeResponseBody(rw, user, "user")

}

// @Summary      Show user list
// @Description  Get user list
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  reqmodel.UserListResponse
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user [get]
func (ho *HandlerObj) GetUserListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	userList, err := crudl.GetUserList(ctx, ho.DBPool)
	if err != nil {
		ho.Logger.Printf("proceed getting user list: %w", err)
		http.Error(rw, "Can't proceed getting user list", http.StatusBadRequest)
		return
	}

	userListResponse := reqmodel.UserListResponse{UserList: userList}
	writeResponseBody(rw, userListResponse, "user list")
}

// @Summary      Update user
// @Description  Update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id   path      string  true  "User ID"
// @Param        request 		body	reqmodel.UserRequest  true  "User creation data"
// @Success      200  {object}  sqlc.UserDatum
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id} [PATCH]
func (ho *HandlerObj) UpdateUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	var userUpdateRequest reqmodel.UserUpdateRequest
	if err := decoder.Decode(&userUpdateRequest); err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	encodedPassword := secret.EncodePassword(*userUpdateRequest.Password)
	userUpdate := sqlc.UpdateUserParams{ID: userID, Name: userUpdateRequest.Name, Login: userUpdateRequest.Login, EncodePassword: encodedPassword}
	user, err := crudl.UpdateUser(ctx, ho.DBPool, userUpdate)
	if err != nil {
		ho.Logger.Printf("proceed update user: %v", err)
		http.Error(rw, "Can't proceed update user", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, user, "user")
}

// @Summary      Create user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request 		body	reqmodel.UserRequest  true  "User creation data"
// @Success      201  {object}  sqlc.UserDatum
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user [post]
func (ho *HandlerObj) CreateUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict parsing

	var userRequest reqmodel.UserRequest
	err := decoder.Decode(&userRequest)
	if err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	encodedPassword := secret.EncodePassword(userRequest.Password)
	userCreate := sqlc.CreateUserParams{Name: userRequest.Name, Login: userRequest.Login, EncodedPassword: encodedPassword}
	user, err := crudl.CreateUser(ctx, ho.DBPool, userCreate)
	if err != nil {
		ho.Logger.Printf("proceed user creation: %v", err)
		http.Error(rw, "Can't create user", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, user, "user")
}

// @Summary  		Delete user
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       user_id   path	string  true  "User ID"
// @Success     204
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /user/{user_id}	[delete]
func (ho *HandlerObj) DeleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	err := crudl.DeleteUser(ctx, ho.DBPool, userID)
	if err != nil {
		ho.Logger.Printf("proceed user deletion: %v", err)
		http.Error(rw, "Can't delete user", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
