package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/llucascr/first_service_go/handler"
	"github.com/llucascr/first_service_go/middleware"
	"github.com/llucascr/first_service_go/config"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)

func main() {
	log.Println("creating a webserver")

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := config.NewPostgresDB(connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to PostgreSQL successfully established!")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	fmt.Println(userService)

	notebookRepository := repository.NewNotebookRepository(db)
	notebookService := service.NewNotebookService(notebookRepository)
	notebookHandler := handler.NewNotebookHandler(notebookService)

	tagRepository := repository.NewTagRepository(db)
	tagService := service.NewTagService(tagRepository)
	tagHandler := handler.NewTagHandler(tagService)

	router := mux.NewRouter()
	router.HandleFunc("/health", handler.Health).Methods("GET")

	notebookHandler.MountNotebookHandler(router)
	tagHandler.MountTagHandler(router)
	loggedRouter := middleware.LoggingMiddleware(router)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))

}
