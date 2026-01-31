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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/swaggo/http-swagger/v2"
)

// @title           					movie_backend_go
// @version         					1.0
// @description     					Basic swagger for current api
// @termsOfService  					http://swagger.io/terms/
// @host      								localhost:8080
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

	// mux := http.NewServeMux()
	handlerObj := handlers.HandlerObj{DBPool: queries, Log: log.Default()}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//auth
	r.Post("/auth/login", handlerObj.LoginHandler)

	// user
	r.Get("/user/{user_id}", handlerObj.GetUserHandler)
	r.Get("/user", handlerObj.GetUserListHandler)
	r.Post("/user", handlerObj.CreateUserHandler)
	r.Patch("/user/{user_id}", handlerObj.UpdateUserHandler)
	r.Delete("/user/{user_id}", handlerObj.DeleteUserHandler)

	r.Get("/user/{user_id}/comment", handlerObj.GetUserCommentListHandler)
	r.Get("/user/{user_id}/rating", handlerObj.GetUserRatingListHandler)
	r.Get("/user/{user_id}/favorite", handlerObj.GetUserFavoriteListHandler)

	// movie
	r.Get("/movie", handlerObj.GetMovieListHandler)
	r.Post("/movie", handlerObj.CreateMovieHandler)
	r.Get("/movie/{movie_id}", handlerObj.GetMovieHandler)
	r.Patch("/movie/{movie_id}", handlerObj.UpdateMovieHandler)
	r.Delete("/movie/{movie_id}", handlerObj.DeleteMovieHandler)

	r.Get("/movie/{movie_id}/comment", handlerObj.GetMovieCommentListHandler)
	r.Get("/movie/{movie_id}/rating", handlerObj.GetMovieRatingListHandler)
	r.Get("/movie/{movie_id}/favorite", handlerObj.GetMovieFavoriteListHandler)

	// rating
	r.Get("/rating", handlerObj.GetRatingHandler)
	r.Post("/rating", handlerObj.CreateRatingHandler)
	r.Patch("/rating", handlerObj.UpdateRatingHandler)
	r.Delete("/rating", handlerObj.DeleteRatingHandler)

	// favorite
	r.Get("/favorite", handlerObj.GetFavoriteHandler)
	r.Post("/favorite", handlerObj.CreateFavoriteHandler)
	r.Delete("/favorite", handlerObj.DeleteFavoriteHandler)

	// System specific commands
	r.Get("/healthcheck", handlers.CheckHealthHandlerCreate(dbPool))
	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
