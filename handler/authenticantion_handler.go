package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/response"
	"github.com/llucascr/first_service_go/service"
)

type AuthenticationHandler struct {
	service *service.AuthenticationService
}

func NewAuthenticationHandler(srv *service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{
		service: srv,
	}
}

func (h *AuthenticationHandler) mountHandler(r *mux.Router) {
	r.HandleFunc("/signup", h.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/signin", h.SignIn).Methods(http.MethodPost)
}

func (h *AuthenticationHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	var request model.SingUpRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, model.ErrInvalidBody)
		return
	}

	// TODO: utils de enviar email
	err := h.service.SingUp(ctx, request)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "Usuário criado com sucesso",
	})
}

func (h *AuthenticationHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	var request model.SingInRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, model.ErrInvalidBody)
		return
	}

	token, err := h.service.SingIn(ctx, request)
	if err != nil {
		response.Error(w, err)
		return
	}

	r.Header.Add("Authorization", token)
	response.JSON(w, http.StatusCreated, map[string]string{
		"token": token,
	})
}
