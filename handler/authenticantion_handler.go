package handler

import (
	// "context"
	// "encoding/json"
	// "fmt"
	// "net/http"

	// "github.com/gorilla/mux"
	// "github.com/llucascr/first_service_go/model"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/service"
)

type UserHandler struct {
	service *service.AuthenticationService
}

func NewUserHandler(srv *service.AuthenticationService) *UserHandler {
	return &UserHandler{
		service: srv,
	}
}

func (h *UserHandler) mountHandler(r *mux.Router) {
	r.HandleFunc("/signup", h.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/signin", h.SignIn).Methods(http.MethodPost)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	var request model.SingUpRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: utils de enviar email
	resp, err := h.service.SingUp(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao criar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	var request model.SingInRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.SingIn(ctx, request)
	if err != nil {
		http.Error(w, "Erro logar o usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
