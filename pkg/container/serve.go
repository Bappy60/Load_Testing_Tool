package container

import (
	"log"
	"net/http"

	"github.com/Bappy60/BookStore_in_Go/pkg/config"
	"github.com/Bappy60/BookStore_in_Go/pkg/connection"
	"github.com/Bappy60/BookStore_in_Go/pkg/controllers"
	"github.com/Bappy60/BookStore_in_Go/pkg/repositories"
	"github.com/Bappy60/BookStore_in_Go/pkg/routes"
	"github.com/Bappy60/BookStore_in_Go/pkg/services"
	"github.com/gorilla/mux"
)

func Serve() {
	config.SetConfig()
	var db = connection.Initialize()
	redisClient:= connection.Redis()

	bookRepo := repositories.BookDBInstance(db)
	bookService := services.BookServiceInstance(bookRepo,redisClient)
	bookController := controllers.BookControllerInstance(bookService)

	services.PopulateBookCacheMap(bookRepo)

	authorRepo := repositories.AuthorDBInstance(db)
	authorService := services.AuthorServiceInstance(authorRepo)
	authorController := controllers.AuthorControllerInstance(authorService)

	log.Println("Database Connected...")
	r := mux.NewRouter()

	// Health Check Endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Registering Routes
	routes.AuthorRoutes(r, authorController)
	routes.BookRoutes(r, bookController)

	// HTTP Server
	log.Println("Server Started...")
	log.Fatal(http.ListenAndServe(":9011", r))
}
