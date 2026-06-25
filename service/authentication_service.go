package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
)

type AuthenticationService struct {
	repository *repository.AuthenticationRepository
}

func NewAuthenticationService(repo *repository.AuthenticationRepository) *AuthenticationService {
	return &AuthenticationService{
		repository: repo,
	}
}

func (srv *AuthenticationService) SingUp(ctx context.Context, dto model.SingUpRequestDTO) error {

	now := time.Now()
	new_user := &model.User{
		UserID:    uuid.New(),
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  dto.Password,
		DeletedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := srv.repository.GetUserByName(ctx, new_user.Name)
	if err == nil {
		return model.ErrUserExists
	}
	if !errors.Is(err, model.ErrNotFound) {
		return fmt.Errorf("get user by name: %w", err)
	}

	_, err = srv.repository.CreateUser(ctx, *new_user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (srv *AuthenticationService) SingIn(ctx context.Context, dto model.SingInRequestDTO) (string, error) {

	user_found, err := srv.repository.GetUserByName(ctx, dto.Name)
	if errors.Is(err, model.ErrNotFound) {
		return "", model.ErrInvalidCredential
	}
	if err != nil {
		return "", fmt.Errorf("get user by name: %w", err)
	}

	if user_found.Password != dto.Password {
		return "", model.ErrInvalidCredential
	}

	token, err := user_found.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("Error Generated Token")
	}

	return token, nil
}

func (srv *AuthenticationService) GetUserByName(ctx context.Context, name string) (*model.User, error) {

	user_found, err := srv.repository.GetUserByName(ctx, name)
	if errors.Is(err, model.ErrNotFound) {
		return nil, model.ErrInvalidCredential
	}
	if err != nil {
		return nil, fmt.Errorf("get user by name: %w", err)
	}

	return user_found, nil
}
