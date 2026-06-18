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

	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
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

	metaContentRepository := repository.NewMetaContentRepository(db)
	metaContentService := service.NewMetaContentService(metaContentRepository)
	metaContentHandler := handler.NewMetaContentHandler(metaContentService)

	nodesContentRepository := repository.NewNodesContentRepository(db)
	nodesContentService := service.NewNodesContentService(nodesContentRepository)
	nodesContentHandler := handler.NewNodesContentHandler(nodesContentService)

	metaTagContentRepository := repository.NewMetaTagContentRepository(db)
	metaTagContentService := service.NewMetaTagContentService(metaTagContentRepository)
	metaTagContentHandler := handler.NewMetaTagContentHandler(metaTagContentService)

	router := mux.NewRouter()
	router.HandleFunc("/health", handler.Health).Methods("GET")

	notebookHandler.MountNotebookHandler(router)
	tagHandler.MountTagHandler(router)
	metaContentHandler.MountMetaContentHandler(router)
	nodesContentHandler.MountNodesContentHandler(router)
	metaTagContentHandler.MountMetaTagContentHandler(router)
	loggedRouter := middleware.LoggingMiddleware(router)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))

}
