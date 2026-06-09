package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/service"
)

type Handler struct {
	service *service.Service
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func New(srv *service.Service) *Handler {
	return &Handler{
		service: srv,
	}
}

func (h *Handler) MountHandler(r *mux.Router) {
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/notebooks", h.Create).Methods("POST")
	r.HandleFunc("/notebooks/list", h.List).Methods("GET")
	r.HandleFunc("/notebooks/{id}", h.Get).Methods("GET")
	r.HandleFunc("/notebooks/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/notebooks/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok", Message: "Service is healthy"}); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
}

// Create
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	ctx := context.TODO()

	var notebookRequestDTO model.NotebookRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&notebookRequestDTO); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %s", err), http.StatusInternalServerError)
		return
	}

	resp, err := h.service.Create(ctx, notebookRequestDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create notebook: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
}

// Update
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	ctx := context.TODO()

	vars := mux.Vars(r)

	var notebookRequestDTO model.NotebookRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&notebookRequestDTO); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %s", err), http.StatusInternalServerError)
		return
	}

	resp, err := h.service.Update(ctx, notebookRequestDTO, vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to updated notebook: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
}

// Get
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.TODO()

	vars := mux.Vars(r)

	notebook, err := h.service.Get(ctx, vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Not Found notebook: %s", err), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(notebook); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
}

// Delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.TODO()

	vars := mux.Vars(r)
	err := h.service.Delete(ctx, vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Not Found notebook: %s", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Sucesso ao deletar documento",
		"documento": vars["id"],
	})
}

// List
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.TODO()

	notebooks, err := h.service.List(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list notebooks: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notebooks); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
}
