package handler

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)

// MountHandler constroi toda a cadeia repository -> service -> handler a partir
// do *sql.DB e delega o registro das rotas para o mountHandler de cada handler.
func MountHandler(r *mux.Router, db *sql.DB) {
	// Health
	r.HandleFunc("/health", health).Methods(http.MethodGet)

	// User
	NewUserHandler(service.NewUserService(repository.NewUserRepository(db))).mountHandler(r)

	// Notebook
	NewNotebookHandler(service.NewNotebookService(repository.NewNotebookRepository(db))).mountHandler(r)

	// Tag
	NewTagHandler(service.NewTagService(repository.NewTagRepository(db))).mountHandler(r)

	// Meta Content
	NewMetaContentHandler(service.NewMetaContentService(repository.NewMetaContentRepository(db))).mountHandler(r)

	// Nodes Content
	NewNodesContentHandler(service.NewNodesContentService(repository.NewNodesContentRepository(db))).mountHandler(r)

	// Meta Tag Content
	NewMetaTagContentHandler(service.NewMetaTagContentService(repository.NewMetaTagContentRepository(db))).mountHandler(r)
}
