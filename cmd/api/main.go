package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/llucascr/first_service_go/config"
	"github.com/llucascr/first_service_go/handler"
	"github.com/llucascr/first_service_go/middleware"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)

func main() {
	log.Println("creating a webserver")

	if err := godotenv.Load(); err != nil {
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

	router := mux.NewRouter()
	handler.MountHandler(router, db)
	loggedRouter := middleware.LoggingMiddleware(router)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))

}
