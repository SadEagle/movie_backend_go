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
	"movie_backend_go/internal/encode"
	"movie_backend_go/internal/handlers/reqmodel"
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
	user, err := crudl.GetUser(ctx, ho.QuerierDB, userID)
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
	userList, err := crudl.GetUserList(ctx, ho.QuerierDB)
	if err != nil {
		ho.Logger.Printf("proceed getting user list: %v", err)
		http.Error(rw, "Can't proceed getting user list", http.StatusBadRequest)
		return
	}

	userListResponse := reqmodel.UserListResponse{UserList: userList}
	writeResponseBody(rw, userListResponse, "user list")
}

// @Summary      Create user
// @Tags         user
// @Accept       multipart/form-data
// @Produce      json
// @Security	 	 OAuth2Password
// @Param        request 		formData	reqmodel.UserCreateRequest  true  "User creation data"
// @Success      201  {object}  sqlc.UserDatum
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user [post]
func (ho *HandlerObj) CreateUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	nameForm := r.FormValue("name")
	if nameForm == "" {
		ho.Logger.Println("Form param `name` is required")
		http.Error(rw, "Form param `name` not found", http.StatusBadRequest)
		return
	}
	loginForm := r.FormValue("login")
	if loginForm == "" {
		ho.Logger.Println("Form param `login` is required")
		http.Error(rw, "Form param `login` not found", http.StatusBadRequest)
		return
	}
	passwordForm := r.FormValue("password")
	if passwordForm == "" {
		ho.Logger.Println("Form param `passoword` is required")
		http.Error(rw, "Form param `password` not found", http.StatusBadRequest)
		return
	}

	encodedPassword := encode.EncodePassword(passwordForm)
	userCreate := sqlc.CreateUserParams{Name: nameForm, Login: loginForm, EncodedPassword: encodedPassword}
	user, err := crudl.CreateUser(ctx, ho.QuerierDB, userCreate)
	if err != nil {
		ho.Logger.Printf("proceed user creation: %v", err)
		http.Error(rw, "Can't create user", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, user, "user")
}

// @Summary      Update user
// @Description  Update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security	 	 OAuth2Password
// @Param        request 		body	reqmodel.UserUpdateRequest  true  "User creation data"
// @Success      200  {object}  sqlc.UserDatum
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/me [patch]
func (ho *HandlerObj) UpdateUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userUpdateRequest reqmodel.UserUpdateRequest
	if err := decoder.Decode(&userUpdateRequest); err != nil && err != io.EOF {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	encodedPassword := encode.EncodePassword(*userUpdateRequest.Password)
	userUpdate := sqlc.UpdateUserParams{ID: userTokenData.UserID, Name: userUpdateRequest.Name, Login: userUpdateRequest.Login, EncodePassword: encodedPassword}
	user, err := crudl.UpdateUser(ctx, ho.QuerierDB, userUpdate)
	if err != nil {
		ho.Logger.Printf("proceed update user: %v", err)
		http.Error(rw, "Can't proceed update user", http.StatusBadRequest)
		return
	}
	writeResponseBody(rw, user, "user")
}

// @Summary  		Delete myself
// @Tags        user
// @Accept      json
// @Produce     json
// @Success     204
// @Failure     401  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /user/me	[delete]
func (ho *HandlerObj) MyselfDeleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	err = crudl.DeleteUser(ctx, ho.QuerierDB, userTokenData.UserID)
	if err != nil {
		ho.Logger.Printf("proceed user deletion: %v", err)
		http.Error(rw, "Can't delete user", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

// @Summary  		Delete user
// @Description	Delete any user as admin
// @Tags        user, admin
// @Accept      json
// @Produce     json
// @Security	 	OAuth2Password
// @Param       user_id   path	string  true  "User ID"
// @Success     204
// @Failure     401  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /user/{user_id}	[delete]
func (ho *HandlerObj) AdminDeleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	var userID pgtype.UUID
	if err := userID.Scan(r.PathValue("user_id")); err != nil {
		ho.Logger.Printf("proceed body request: %v", err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	// Extract token
	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}
	// Verify admin
	if !userTokenData.IsAdmin {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	err = crudl.DeleteUser(ctx, ho.QuerierDB, userID)
	if err != nil {
		ho.Logger.Printf("proceed user deletion: %v", err)
		http.Error(rw, "Can't delete user", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
