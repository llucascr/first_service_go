package main

import (
	"log"
	"net/http"

	kivik "github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb" // The CouchDB driver
	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/handler"
	"github.com/llucascr/first_service_go/middleware"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)

func main() {
	log.Println("creating a webserver")

	client, err := kivik.New("couch", "http://admin:pass@localhost:5984/")
	if err != nil {
		log.Fatalf("Failed to create client: %s", err)
	}

	db := client.DB("notebooks")
	if err := db.Err(); err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)
	srv := service.NewService(repo)
	handler := handler.New(srv)
	
	router := mux.NewRouter()

	handler.MountHandler(router)

	loggedRouter := middleware.LoggingMiddleware(router)

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))

}