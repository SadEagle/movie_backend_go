package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"movie_backend_go/crudl"
	"movie_backend_go/models"
	"time"

	"net/http"
)

// @Summary      Show user
// @Description  Get user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.UserData
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{id} [get]
func (ho *HandlerObj) GetUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()

	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	user, err := crudl.GetUserByIDDB(ctx, ho.DB, userID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't get user id: %s\n", userID), http.StatusNotFound)
		return
	}
	writeResponseBody(rw, user, "user")

}

// @Summary      Show user list
// @Description  Get user list
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.UserDataList
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user [get]
func (ho *HandlerObj) GetUserListHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	userList, err := crudl.GetUserListDB(ctx, ho.DB)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't get user list", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	writeResponseBody(rw, userList, "user list")
}

// @Summary      Update user
// @Description  Update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        request 		body	models.UpdateUserDataRequest  true  "User creation data"
// @Success      200  {object}  models.UserData
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{id} [PATCH]
func (ho *HandlerObj) UpdateUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict parsing

	var updateUserdata models.UpdateUserRequest
	err := decoder.Decode(&updateUserdata)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	fmt.Println(updateUserdata)
	user, err := crudl.UpdateUserDB(ctx, ho.DB, updateUserdata, userID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't update user", http.StatusNotFound)
		return
	}

	writeResponseBody(rw, user, "user")
}

// @Summary      Create user
// @Description  Create user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request 		body	models.CreateUserDataRequest  true  "User creation data"
// @Success      201  {object}  models.UserData
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user [post]
func (ho *HandlerObj) CreateUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict parsing

	var createUserData models.CreateUserRequest
	err := decoder.Decode(&createUserData)
	if err != nil && err != io.EOF {
		ho.Log.Println(err)
		http.Error(rw, "Can't proceed body request", http.StatusBadRequest)
		return
	}

	user, err := crudl.CreateUserDB(ctx, ho.DB, createUserData)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Can't create user", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	writeResponseBody(rw, user, "user")
}

// @Description  Delete user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  {object}  models.UserData
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{id} [delete]
func (ho *HandlerObj) DeleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), 5*time.Minute)
	defer close()

	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, "Requested user id should contain uuid style", http.StatusBadRequest)
		return
	}

	err = crudl.DeleteUserDB(ctx, ho.DB, userID)
	if err != nil {
		ho.Log.Println(err)
		http.Error(rw, fmt.Sprintf("Can't delete user id: %s", userID), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
