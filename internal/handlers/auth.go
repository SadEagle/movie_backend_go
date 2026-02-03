package handlers

import (
	"bytes"
	"context"
	_ "golang.org/x/oauth2"
	"movie_backend_go/internal/auth"
	"movie_backend_go/internal/crudl"
	"movie_backend_go/internal/encode"
	"net/http"
	"time"
)

var EXPIRE_TIME = 24 * time.Hour

// @Summary      Auth
// @Description  Get movie by id
// @Tags         auth
// @Accept       multipart/form-data
// @Produce      json
// @Param        username		formData	string  true  "Login"
// @Param        password		formData	string  true  "Password"
// @Success      200  {object}  oauth2.Token
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/login [post]
func (ho *HandlerObj) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	loginStr := r.FormValue("username")
	if loginStr == "" {
		ho.Logger.Println("Form param `username` is required")
		http.Error(rw, "Form param `username` not found", http.StatusBadRequest)
		return
	}

	passwordStr := r.FormValue("password")
	if passwordStr == "" {
		ho.Logger.Println("Form param `password` is required")
		http.Error(rw, "Form param `password` not found", http.StatusBadRequest)
		return
	}

	user, err := crudl.GetUserByLogin(ctx, ho.QuerierDB, loginStr)
	if err != nil {
		ho.Logger.Printf("get user by login from db: %v", err)
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}
	if !bytes.Equal(encode.EncodePassword(passwordStr), user.EncodedPassword) {
		ho.Logger.Printf("Wrong password")
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	userTokenData := auth.UserTokenData{UserID: user.ID, IsAdmin: user.IsAdmin}
	oauthToken, err := auth.OauthTokenGenerate(userTokenData)
	if err != nil {
		ho.Logger.Printf("generate token: %v", err)
		http.Error(rw, "Can't generate user token", http.StatusInternalServerError)
		return
	}
	writeResponseBody(rw, oauthToken, "oauth token")
}
