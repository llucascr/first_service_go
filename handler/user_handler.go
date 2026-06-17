package handler

import (
	// "context"
	// "encoding/json"
	// "fmt"
	// "net/http"

	// "github.com/gorilla/mux"
	// "github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{
		service: srv,
	}
}