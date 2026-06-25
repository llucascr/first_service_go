package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/llucascr/first_service_go/model"
)

type AuthenticationRepository struct {
	database *sql.DB
}

//go:embed queries/user_get_by_name.sql
var getUserByNameQuery string

//go:embed queries/user_create.sql
var createUserQuery string

func NewAuthenticationRepository(db *sql.DB) *AuthenticationRepository {
	return &AuthenticationRepository{
		database: db,
	}
}

func (r *AuthenticationRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	_, err := r.database.ExecContext(
		ctx,
		createUserQuery,
		user.UserID,
		user.Name,
		user.Email,
		user.Password,
		user.DeletedAt,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthenticationRepository) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.database.QueryRowContext(ctx, getUserByNameQuery, name).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
