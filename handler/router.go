package handler

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/llucascr/first_service_go/middleware"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)

func MountHandler(r *mux.Router, db *sql.DB) {
	r.Use(middleware.LoggingMiddleware)
	r.HandleFunc("/health", health).Methods(http.MethodGet)

	auth := r.PathPrefix("/auth").Subrouter()
	// Auth (User)
	authService := service.NewAuthenticationService(repository.NewAuthenticationRepository(db))
	NewAuthenticationHandler(authService).mountHandler(auth)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.Identify(authService))

	// Notebook
	NewNotebookHandler(service.NewNotebookService(repository.NewNotebookRepository(db))).mountHandler(api)

	// Tag
	NewTagHandler(service.NewTagService(repository.NewTagRepository(db))).mountHandler(api)

	// Meta Content
	NewMetaContentHandler(service.NewMetaContentService(repository.NewMetaContentRepository(db))).mountHandler(api)

	// Nodes Content
	NewNodesContentHandler(service.NewNodesContentService(repository.NewNodesContentRepository(db))).mountHandler(api)

	// Meta Tag Content
	NewMetaTagContentHandler(service.NewMetaTagContentService(repository.NewMetaTagContentRepository(db))).mountHandler(api)
}
