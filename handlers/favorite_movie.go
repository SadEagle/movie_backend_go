package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"movie_backend_go/db/crudl"
	"movie_backend_go/models"
	"net/http"
)

// func writeFavMovieResponseBody(rw http.ResponseWriter, favMovie models.FavoriteMovie) {
// 	responseFavMovie := favMovie.ToResponse()
// 	responseFavMovieByte, err := json.Marshal(responseFavMovie)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(rw, "Can't convert favorite movie data to response json", 500)
// 	}
//
// 	_, err = rw.Write(responseFavMovieByte)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(rw, "Can't send favorite movie data", 500)
// 	}
// }

// @Description  Get user favorite_movie list
// @Tags         favorite_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Success      200  {object}  models.FavoriteMovieListResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/favorite_movie [get]
func GetFavoriteMovieListHandlerMake(db *sql.DB) http.HandlerFunc {
	AddFavoriteMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		user_id := r.PathValue("user_id")
		favMovieList, err := crudl.GetFavoriteMovieListDB(db, user_id)
		if err != nil {
			log.Println(err)
			http.Error(rw, fmt.Sprintf("Can't get favorite_movie list of user: %s", user_id), 404)
			return
		}
		writeResponseBody(rw, favMovieList, "favorite_move")
	}
	return AddFavoriteMovieHandler
}

// @Description  Add favorite movie
// @Tags         favorite_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      200  {object}  models.FavoriteMovieResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/movie/{movie_id} [post]
func AddFavoriteMovieHandlerMake(db *sql.DB) http.HandlerFunc {
	AddFavoriteMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		user_id := r.PathValue("user_id")
		movie_id := r.PathValue("movie_id")
		favMovie, err := crudl.AddFavoriteMovieDB(db, user_id, movie_id)
		if err != nil {
			log.Println(err)
			http.Error(rw, fmt.Sprintf("Can't add favorite movie user_id: %s, movie_id: %s", user_id, user_id), 404)
			return
		}
		rw.WriteHeader(201) // 201 - Create
		writeFavMovieResponseBody(rw, favMovie)
	}
	return AddFavoriteMovieHandler
}

// @Description  Delete favorite movie
// @Tags         favorite_movie
// @Accept       json
// @Produce      json
// @Param        user_id   	path      string  true  "User ID"
// @Param        movie_id   path      string  true  "Movie ID"
// @Success      204  {object}  models.FavoriteMovieResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{user_id}/{movie_id} [delete]
func DeleteFavoriteMovieHandlerMake(db *sql.DB) http.HandlerFunc {
	AddFavoriteMovieHandler := func(rw http.ResponseWriter, r *http.Request) {
		user_id := r.PathValue("user_id")
		movie_id := r.PathValue("movie_id")
		err := crudl.DeleteFavoriteMovieDB(db, user_id, movie_id)
		if err != nil {
			log.Println(err)
			http.Error(rw, fmt.Sprintf("Can't delete favorite movie user_id: %s, movie_id: %s", user_id, user_id), 404)
			return
		}
		rw.WriteHeader(204) // 204 - No body
	}
	return AddFavoriteMovieHandler
}
