package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/response"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthService interface {
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}

func Identify(auth AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("Authorization")
			if token == "" {
				response.Error(w, model.ErrInvalidCredential)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")

			name, err := model.ValidateToken(token)
			if err != nil {
				response.Error(w, model.ErrInvalidCredential)
				return
			}

			user, err := auth.GetUserByName(r.Context(), name)
			if err != nil {
				response.Error(w, model.ErrInvalidCredential)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}