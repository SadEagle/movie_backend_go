package handlers

import (
	"bytes"
	"context"
	"movie_backend_go/crudl"
	"movie_backend_go/secret"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var EXPIRE_TIME = 24 * time.Hour

// @Summary      Auth
// @Description  Get movie by id
// @Tags         auth
// @Accept       multipart/form-data
// @Produce      json
// @Param        login	formData	string  true  "Login"
// @Param        password	formData	string  true  "Password"
// @Success      200  {object}  oauth2.Token
// @Failure      404  {object}	map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/login [post]
func (ho *HandlerObj) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	loginStr := r.FormValue("login")
	if loginStr == "" {
		ho.Log.Println("Form param `login` is required")
		http.Error(rw, "Form param `login` not found", http.StatusBadRequest)
		return
	}

	passwordStr := r.FormValue("password")
	if passwordStr == "" {
		ho.Log.Println("Form param `password` is required")
		http.Error(rw, "Form param `password` not found", http.StatusBadRequest)
		return
	}

	user, err := crudl.GetUserByLogin(ctx, ho.DBPool, loginStr)
	if err != nil {
		ho.Log.Printf("get user by login from db: %v", err)
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	if !bytes.Equal(secret.EncodePassword(passwordStr), user.EncodedPassword) {
		ho.Log.Printf("Wrong formData password: %v", err)
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	experify := time.Now().Add(EXPIRE_TIME)
	expiresIn := int64(EXPIRE_TIME.Seconds())
	userTokenData := secret.UserTokenData{UserID: user.ID, IsAdmin: user.IsAdmin}

	accessToken, err := secret.TokenGenerate(userTokenData, experify)
	if err != nil {
		ho.Log.Printf("generate token: %v", err)
		http.Error(rw, "Can't generate user token", http.StatusInternalServerError)
	}

	token := oauth2.Token{AccessToken: accessToken, TokenType: "Bearer", Expiry: experify, ExpiresIn: expiresIn}
	writeResponseBody(rw, token, "token")
}
