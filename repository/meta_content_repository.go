package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
)

type MetaContentRepository struct {
	database *sql.DB
}

//go:embed queries/meta_content_create.sql
var createMetaContentQuery string

//go:embed queries/meta_content_list_from_notebook.sql
var listMetaContentFromNotebookQuery string

//go:embed queries/meta_content_get_by_id.sql
var getMetaContentByIDQuery string

//go:embed queries/meta_content_update.sql
var updateMetaContentQuery string

//go:embed queries/meta_content_delete.sql
var deleteMetaContentQuery string

func NewMetaContentRepository(db *sql.DB) *MetaContentRepository {
	return &MetaContentRepository{
		database: db,
	}
}

func (r *MetaContentRepository) Create(ctx context.Context, content model.MetaContent) (*model.MetaContent, error) {
	_, err := r.database.ExecContext(
		ctx,
		createMetaContentQuery,
		content.ContentID,
		content.NotebookID,
		content.UserID,
		content.Icon,
		content.Name,
		content.DeletedAt,
		content.CreatedAt,
		content.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (r *MetaContentRepository) ListMetaContentsFromNotebook(ctx context.Context, notebookID uuid.UUID) ([]*model.MetaContent, error) {
	rows, err := r.database.QueryContext(ctx, listMetaContentFromNotebookQuery, notebookID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var contents []*model.MetaContent

	for rows.Next() {
		var content model.MetaContent
		if err := rows.Scan(
			&content.ContentID,
			&content.NotebookID,
			&content.UserID,
			&content.Icon,
			&content.Name,
			&content.CreatedAt,
			&content.UpdatedAt,
		); err != nil {
			return nil, err
		}
		contents = append(contents, &content)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}

func (r *MetaContentRepository) GetMetaContentByID(ctx context.Context, contentID uuid.UUID) (*model.MetaContent, error) {
	var content model.MetaContent
	err := r.database.QueryRowContext(ctx, getMetaContentByIDQuery, contentID).Scan(
		&content.ContentID,
		&content.NotebookID,
		&content.UserID,
		&content.Icon,
		&content.Name,
		&content.DeletedAt,
		&content.CreatedAt,
		&content.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (r *MetaContentRepository) Update(ctx context.Context, content model.MetaContent) (*model.MetaContent, error) {
	var updated model.MetaContent
	err := r.database.QueryRowContext(
		ctx,
		updateMetaContentQuery,
		content.Icon,
		content.Name,
		content.UpdatedAt,
		content.ContentID,
	).Scan(
		&updated.ContentID,
		&updated.NotebookID,
		&updated.UserID,
		&updated.Icon,
		&updated.Name,
		&updated.DeletedAt,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *MetaContentRepository) Delete(ctx context.Context, contentID uuid.UUID, deletedAt time.Time) error {
	_, err := r.database.ExecContext(ctx, deleteMetaContentQuery, deletedAt, contentID)
	if err != nil {
		return err
	}

	return nil
}
