package main

import (
	"fmt"
	"log"
	"movie_backend_go/db"
	"movie_backend_go/db/sqlc"
	_ "movie_backend_go/docs"
	"movie_backend_go/handlers"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/swaggo/http-swagger/v2"
)

// @title           movie_backend_go
// @version         1.0
// @description     Basic swagger for current api
// @termsOfService  http://swagger.io/terms/
// @host      localhost:8080
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalln(fmt.Errorf("parsing port value: %w", err))
	}
	text, err := os.ReadFile(os.Getenv("DB_PASSWORD_FILE"))
	if err != nil {
		log.Fatalln(fmt.Errorf("reading secret file: %w", err))
	}
	c := db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: string(text),
		SSLMode:  "disable",
	}

	dbPool, err := db.InitDB(c)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbPool.Close()
	queries := sqlc.New(dbPool)

	ping_db_check := func() {
		for {
			// err = db.Ping()
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Minute * 5)
		}
	}
	go ping_db_check()

	mux := http.NewServeMux()
	handlerObj := handlers.HandlerObj{DBPool: queries, Log: log.Default()}

	// user
	mux.HandleFunc("GET /user/{user_id}", handlerObj.GetUserHandler)
	mux.HandleFunc("GET /user", handlerObj.GetUserListHandler)
	mux.HandleFunc("POST /user", handlerObj.CreateUserHandler)
	mux.HandleFunc("PATCH /user/{user_id}", handlerObj.UpdateUserHandler)
	mux.HandleFunc("DELETE /user/{user_id}", handlerObj.DeleteUserHandler)
	// rating
	mux.HandleFunc("GET /user/{user_id}/rating", handlerObj.GetRatedMovieListHandler)
	mux.HandleFunc("POST /user/{user_id}/rating", handlerObj.CreateRatedMovieHandler)
	mux.HandleFunc("PATCH /user/{user_id}/rating", handlerObj.UpdateRatedMovieHandler)
	mux.HandleFunc("DELETE /user/{user_id}/rating/{movie_id}", handlerObj.DeleteRatedMovieHandler)
	// fav
	mux.HandleFunc("GET /user/{user_id}/favorite_movie", handlerObj.GetFavoriteMovieListHandler)
	mux.HandleFunc("POST /user/{user_id}/favorite_movie", handlerObj.CreateMovieFavoriteHandler)
	mux.HandleFunc("PATCH /user/{user_id}/favorite_movie", handlerObj.UpdateMovieHandler)
	mux.HandleFunc("DELETE /user/{user_id}/favorite_movie/{movie_id}", handlerObj.DeleteFavoriteMovieHandler)
	// comment
	mux.HandleFunc("GET /user/{user_id}/comment", handlerObj.GetMovieCommentListHandler)
	mux.HandleFunc("POST /user/{user_id}/comment", handlerObj.CreateMovieCommentHandler)
	mux.HandleFunc("PATCH /user/{user_id}/comment", handlerObj.UpdateMovieCommentHandler)
	mux.HandleFunc("DELETE /user/{user_id}/comment", handlerObj.DeleteMovieCommentHandler)
	// movie
	mux.HandleFunc("GET /movie/{movie_id}", handlerObj.GetMovieHandler)
	mux.HandleFunc("GET /movie", handlerObj.GetMovieListHandler)
	mux.HandleFunc("POST /movie", handlerObj.CreateMovieHandler)
	mux.HandleFunc("PATCH /movie/{movie_id}", handlerObj.UpdateMovieHandler)
	mux.HandleFunc("DELETE /movie/{movie_id}", handlerObj.DeleteMovieHandler)
	// comment
	mux.HandleFunc("GET /movie/{movie_id}/comment", handlerObj.GetMovieCommentListHandler)

	// System specific commands
	mux.Handle("GET /healthcheck", handlers.CheckHealthHandlerCreate(dbPool))
	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
