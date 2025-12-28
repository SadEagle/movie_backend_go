package main

import (
	"github.com/swaggo/http-swagger/v2"
	"log"
	"movie_backend_go/db"
	_ "movie_backend_go/docs"
	"movie_backend_go/handlers"
	"net/http"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{id}", handlers.GetUserHandlerMake(db))
	mux.HandleFunc("POST /user", handlers.CreateUserHandler(db))
	mux.HandleFunc("DELETE /user/{id}", handlers.DeleteUserHandler(db))

	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	// mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	http.ServeFile(w, r, "./docs/swagger.json")
	// })

	log.Fatal(http.ListenAndServe(":8080", mux))
}
