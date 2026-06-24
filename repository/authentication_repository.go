package repository

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/llucascr/first_service_go/model"
)

type UserRepository struct {
	database *sql.DB
}

//go:embed queries/user_get_by_name.sql
var getUserByNameQuery string

//go:embed queries/user_create.sql
var createUserQuery string

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func(r *UserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	_, err := r.database.Exec(
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

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.database.QueryRowContext(ctx, getUserByNameQuery, name).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
