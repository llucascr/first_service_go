package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/utils"
)

type AuthenticationService struct {
	repository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *AuthenticationService {
	return &AuthenticationService{
		repository: repo,
	}
}

func (srv *AuthenticationService) SingUp(ctx context.Context, dto model.SingUpRequestDTO) (*model.UserAccess, error) {

	now := time.Now()
	new_user := &model.User{
		UserID: uuid.New(),
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		DeletedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := srv.repository.GetUserByName(ctx, new_user.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Internal Error")
	}

	if err == nil {
		log.Println("[Authentication][Service] Usuário já existe")
		return nil, errors.New("Invalid Credential")
	}

	saved_user, err := srv.repository.CreateUser(ctx, *new_user)
	if err != nil {
		return nil, err
	}

	return &model.UserAccess{
		User: *saved_user,
		Token: utils.ToBase64(saved_user.Name),
	}, nil
}

func (srv *AuthenticationService) SingIn(ctx context.Context, dto model.SingInRequestDTO) (*model.UserAccess, error) {

	user_found, err := srv.repository.GetUserByName(ctx, dto.Name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Invalid Credential")
	}

	if user_found.Password != dto.Password {
		return nil, errors.New("Invalid Credential")
	}

	return &model.UserAccess{
		User: *user_found,
		Token: utils.ToBase64(user_found.Name),
	}, nil
}
