package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
)

// ErrMetaTagContentAlreadyExists e retornado quando ja existe uma associacao
// (content_id, tag_id) ativa, ou seja, nao soft-deletada.
var ErrMetaTagContentAlreadyExists = errors.New("meta tag content already exists")

type MetaTagContentRepository struct {
	database *sql.DB
}

//go:embed queries/meta_tag_content_create.sql
var createMetaTagContentQuery string

//go:embed queries/meta_tag_content_get.sql
var getMetaTagContentQuery string

//go:embed queries/meta_tag_content_list_from_content.sql
var listMetaTagContentFromContentQuery string

//go:embed queries/meta_tag_content_list_from_tag.sql
var listMetaTagContentFromTagQuery string

//go:embed queries/meta_tag_content_list_from_notebook.sql
var listMetaTagContentFromNotebookQuery string

//go:embed queries/meta_tag_content_update.sql
var updateMetaTagContentQuery string

//go:embed queries/meta_tag_content_delete.sql
var deleteMetaTagContentQuery string

func NewMetaTagContentRepository(db *sql.DB) *MetaTagContentRepository {
	return &MetaTagContentRepository{
		database: db,
	}
}

func (r *MetaTagContentRepository) Create(ctx context.Context, metaTag model.MetaTagContent) (*model.MetaTagContent, error) {
	var saved model.MetaTagContent
	err := r.database.QueryRowContext(
		ctx,
		createMetaTagContentQuery,
		metaTag.ContentID,
		metaTag.TagID,
		metaTag.NotebookID,
		metaTag.UserID,
		metaTag.DeletedAt,
		metaTag.CreatedAt,
		metaTag.UpdatedAt,
	).Scan(
		&saved.ContentID,
		&saved.TagID,
		&saved.NotebookID,
		&saved.UserID,
		&saved.DeletedAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	// Sem linha retornada significa que houve conflito de PK numa associacao ja
	// ativa (o DO UPDATE so reativa linhas soft-deletadas), entao ja existe.
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMetaTagContentAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return &saved, nil
}

func (r *MetaTagContentRepository) ListMetaTagContentsFromContent(ctx context.Context, contentID uuid.UUID) ([]*model.MetaTagContent, error) {
	rows, err := r.database.QueryContext(ctx, listMetaTagContentFromContentQuery, contentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var metaTags []*model.MetaTagContent

	for rows.Next() {
		var metaTag model.MetaTagContent
		if err := rows.Scan(
			&metaTag.ContentID,
			&metaTag.TagID,
			&metaTag.NotebookID,
			&metaTag.UserID,
			&metaTag.CreatedAt,
			&metaTag.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metaTags = append(metaTags, &metaTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metaTags, nil
}

func (r *MetaTagContentRepository) ListMetaTagContentsFromTag(ctx context.Context, tagID uuid.UUID) ([]*model.MetaTagContent, error) {
	rows, err := r.database.QueryContext(ctx, listMetaTagContentFromTagQuery, tagID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var metaTags []*model.MetaTagContent

	for rows.Next() {
		var metaTag model.MetaTagContent
		if err := rows.Scan(
			&metaTag.ContentID,
			&metaTag.TagID,
			&metaTag.NotebookID,
			&metaTag.UserID,
			&metaTag.CreatedAt,
			&metaTag.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metaTags = append(metaTags, &metaTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metaTags, nil
}

func (r *MetaTagContentRepository) ListMetaTagContentsFromNotebook(ctx context.Context, notebookID uuid.UUID) ([]*model.MetaTagContent, error) {
	rows, err := r.database.QueryContext(ctx, listMetaTagContentFromNotebookQuery, notebookID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var metaTags []*model.MetaTagContent

	for rows.Next() {
		var metaTag model.MetaTagContent
		if err := rows.Scan(
			&metaTag.ContentID,
			&metaTag.TagID,
			&metaTag.NotebookID,
			&metaTag.UserID,
			&metaTag.CreatedAt,
			&metaTag.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metaTags = append(metaTags, &metaTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metaTags, nil
}

func (r *MetaTagContentRepository) GetMetaTagContent(ctx context.Context, contentID uuid.UUID, tagID uuid.UUID) (*model.MetaTagContent, error) {
	var metaTag model.MetaTagContent
	err := r.database.QueryRowContext(ctx, getMetaTagContentQuery, contentID, tagID).Scan(
		&metaTag.ContentID,
		&metaTag.TagID,
		&metaTag.NotebookID,
		&metaTag.UserID,
		&metaTag.DeletedAt,
		&metaTag.CreatedAt,
		&metaTag.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &metaTag, nil
}

func (r *MetaTagContentRepository) Update(ctx context.Context, metaTag model.MetaTagContent) (*model.MetaTagContent, error) {
	var updated model.MetaTagContent
	err := r.database.QueryRowContext(
		ctx,
		updateMetaTagContentQuery,
		metaTag.NotebookID,
		metaTag.UpdatedAt,
		metaTag.ContentID,
		metaTag.TagID,
	).Scan(
		&updated.ContentID,
		&updated.TagID,
		&updated.NotebookID,
		&updated.UserID,
		&updated.DeletedAt,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *MetaTagContentRepository) Delete(ctx context.Context, contentID uuid.UUID, tagID uuid.UUID, deletedAt time.Time) error {
	_, err := r.database.ExecContext(ctx, deleteMetaTagContentQuery, deletedAt, contentID, tagID)
	if err != nil {
		return err
	}

	return nil
}
