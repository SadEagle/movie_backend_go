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
	Host:     "dev-db",
	Port:     5432,
	Database: "movie_server",
	User:     "movie_manager",
	Password: "dev_passwd",
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
	handlerObj := handlers.HandlerObj{DB: db, Log: *log.Default()}

	// user
	mux.HandleFunc("GET /user/{id}", handlerObj.GetUserHandler)
	mux.HandleFunc("GET /user", handlerObj.GetUserListHandler)
	mux.HandleFunc("POST /user", handlerObj.CreateUserHandler)
	mux.HandleFunc("PATCH /user/{id}", handlerObj.UpdateUserHandler)
	mux.HandleFunc("DELETE /user/{id}", handlerObj.DeleteUserHandler)
	// movie
	mux.HandleFunc("GET /movie/{id}", handlerObj.GetMovieHandler)
	mux.HandleFunc("GET /movie", handlerObj.GetMovieListHandler)
	mux.HandleFunc("POST /movie", handlerObj.CreateMovieHandler)
	mux.HandleFunc("PATCH /movie/{id}", handlerObj.UpdateMovieHandler)
	mux.HandleFunc("DELETE /movie/{id}", handlerObj.DeleteMovieHandler)
	// favorite_movie
	mux.HandleFunc("GET /user/{user_id}/favorite_movie", handlerObj.GetFavoriteMovieListHandler)
	mux.HandleFunc("POST /user/{user_id}/favorite_movie/{movie_id}", handlerObj.AddFavoriteMovieHandler)
	mux.HandleFunc("DELETE /user/{user_id}/favorite_movie/{movie_id}", handlerObj.DeleteFavoriteMovieHandler)
	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
