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

	var request model.SingUpRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, model.ErrInvalidBody)
		return
	}

	// TODO: utils de enviar email
	resp, err := h.service.SingUp(ctx, request)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	var request model.SingInRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, model.ErrInvalidBody)
		return
	}

	resp, err := h.service.SingIn(ctx, request)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
