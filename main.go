package main

import (
	"github.com/swaggo/http-swagger/v2"
	"log"
	"movie_backend_go/db"
	_ "movie_backend_go/docs"
	"movie_backend_go/handlers"
	"net/http"
	"time"
)

var c = db.Config{
	Host:     "localhost",
	Port:     5432,
	Database: "movie_server",
	User:     "movie_manager",
	Password: "123",
	SSLMode:  "disable",
}

// @title           movie_backend_go
// @version         1.0
// @description     Basic swagger for current api
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	db, err := db.InitDB(c)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ping_db_check := func() {
		for {
			err = db.Ping()
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Minute * 5)
		}
	}
	go ping_db_check()

	mux := http.NewServeMux()

	// user
	updateUserHandler := handlers.UpdateUserHandlerMake(db)

	mux.HandleFunc("GET /user/{id}", handlers.GetUserHandlerMake(db))
	mux.HandleFunc("POST /user", handlers.CreateUserHandlerMake(db))
	mux.HandleFunc("PATCH /user/{id}", updateUserHandler)
	mux.HandleFunc("DELETE /user/{id}", handlers.DeleteUserHandler(db))
	// movie
	mux.HandleFunc("GET /movie/{id}", handlers.GetMovieHandlerMake(db))
	mux.HandleFunc("POST /movie", handlers.CreateMovieHandlerMake(db))
	updateMovieHandler := handlers.UpdateMovieHandlerMake(db)
	mux.HandleFunc("PATCH /movie/{id}", updateMovieHandler)
	mux.HandleFunc("DELETE /movie/{id}", handlers.DeleteMovieHandlerMake(db))
	// favorite_movie
	mux.HandleFunc("GET /user/{user_id}/favorite_movie", handlers.GetFavoriteMovieListHandlerMake(db))
	mux.HandleFunc("POST /user/{user_id}/favorite_movie/{movie_id}", handlers.AddFavoriteMovieHandlerMake(db))
	mux.HandleFunc("DELETE /user/{user_id}/favorite_movie/{movie_id}", handlers.DeleteFavoriteMovieHandlerMake(db))
	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
