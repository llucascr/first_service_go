package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/utils"
)

type AuthenticationService struct {
	repository *repository.AuthenticationRepository
}

func NewAuthenticationService(repo *repository.AuthenticationRepository) *AuthenticationService {
	return &AuthenticationService{
		repository: repo,
	}
}

func (srv *AuthenticationService) SingUp(ctx context.Context, dto model.SingUpRequestDTO) (*model.UserAccess, error) {

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
		return nil, model.ErrUserExists
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("get user by name: %w", err)
	}

	saved_user, err := srv.repository.CreateUser(ctx, *new_user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &model.UserAccess{
		User:  *saved_user,
		Token: utils.ToBase64(saved_user.Name),
	}, nil
}

func (srv *AuthenticationService) SingIn(ctx context.Context, dto model.SingInRequestDTO) (*model.UserAccess, error) {

	user_found, err := srv.repository.GetUserByName(ctx, dto.Name)
	if errors.Is(err, model.ErrNotFound) {
		return nil, model.ErrInvalidCredential
	}
	if err != nil {
		return nil, fmt.Errorf("get user by name: %w", err)
	}

	if user_found.Password != dto.Password {
		return nil, model.ErrInvalidCredential
	}

	return &model.UserAccess{
		User:  *user_found,
		Token: utils.ToBase64(user_found.Name),
	}, nil
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
