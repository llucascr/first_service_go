package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
)

type NotebookRepository struct {
	database *sql.DB
}

//go:embed queries/notebook_create.sql
var createNotebookQuery string

//go:embed queries/notebook_list_from_user.sql
var listNotebookFromUserQuery string

//go:embed queries/notebook_get_by_id.sql
var getNotebookByIDQuery string

//go:embed queries/notebook_update.sql
var updateNotebookQuery string

//go:embed queries/notebook_delete.sql
var deleteNotebookQuery string

func NewNotebookRepository(db *sql.DB) *NotebookRepository {
	return &NotebookRepository{
		database: db,
	}
}

func (r *NotebookRepository) Create(ctx context.Context, notebook model.Notebook) (*model.Notebook, error) {
	_, err := r.database.Exec(
		createNotebookQuery, 
		notebook.NotebookID, 
		notebook.UserID, 
		notebook.Icon, 
		notebook.Name, 
		notebook.Image, 
		notebook.Description,
		notebook.DeletedAt, 
		notebook.CreatedAt, 
		notebook.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &notebook, nil
}

func (r *NotebookRepository) ListNotebooksFromUser(ctx context.Context, userID uuid.UUID) ([]*model.Notebook, error) {
	rows, err := r.database.QueryContext(ctx, listNotebookFromUserQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notebooks []*model.Notebook

	for rows.Next() {
		var notebook model.Notebook
		if err := rows.Scan(
			&notebook.NotebookID,
			&notebook.UserID,
			&notebook.Icon,
			&notebook.Name,
			&notebook.Image,
			&notebook.Description,
			&notebook.CreatedAt,
			&notebook.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notebooks = append(notebooks, &notebook)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notebooks, nil
}

func (r *NotebookRepository) GetNotebookByID(ctx context.Context, notebook_id uuid.UUID) (*model.Notebook, error) {
	var notebook model.Notebook
	err := r.database.QueryRowContext(ctx, getNotebookByIDQuery, notebook_id).Scan(
		&notebook.NotebookID,
		&notebook.UserID,
		&notebook.Name,
		&notebook.Description,
		&notebook.Icon,
		&notebook.Image,
		&notebook.DeletedAt,
		&notebook.CreatedAt,
		&notebook.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &notebook, nil
}

func (r *NotebookRepository) Update(ctx context.Context, notebook model.Notebook) (*model.Notebook, error) {
	var updated model.Notebook
	err := r.database.QueryRowContext(
		ctx,
		updateNotebookQuery,
		notebook.Icon,
		notebook.Name,
		notebook.Image,
		notebook.Description,
		notebook.UpdatedAt,
		notebook.NotebookID,
	).Scan(
		&updated.NotebookID,
		&updated.UserID,
		&updated.Icon,
		&updated.Name,
		&updated.Image,
		&updated.Description,
		&updated.DeletedAt,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *NotebookRepository) Delete(ctx context.Context, notebookID uuid.UUID, deletedAt time.Time) error {
	_, err := r.database.ExecContext(ctx, deleteNotebookQuery, deletedAt, notebookID)
	if err != nil {
		return err
	}

	return nil
}