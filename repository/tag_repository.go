package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
)

type TagRepository struct {
	database *sql.DB
}

//go:embed queries/tag_create.sql
var createTagQuery string

//go:embed queries/tag_list_from_user.sql
var listTagFromUserQuery string

//go:embed queries/tag_get_by_id.sql
var getTagByIDQuery string

//go:embed queries/tag_update.sql
var updateTagQuery string

//go:embed queries/tag_delete.sql
var deleteTagQuery string

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{
		database: db,
	}
}

func (r *TagRepository) Create(ctx context.Context, tag model.Tag) (*model.Tag, error) {
	_, err := r.database.ExecContext(
		ctx,
		createTagQuery,
		tag.TagID,
		tag.Name,
		tag.Color,
		tag.UserID,
		tag.DeletedAt,
		tag.CreatedAt,
		tag.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *TagRepository) ListTagsFromUser(ctx context.Context, userID uuid.UUID) ([]*model.Tag, error) {
	rows, err := r.database.QueryContext(ctx, listTagFromUserQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []*model.Tag

	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(
			&tag.TagID,
			&tag.Name,
			&tag.Color,
			&tag.UserID,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *TagRepository) GetTagByID(ctx context.Context, tagID uuid.UUID) (*model.Tag, error) {
	var tag model.Tag
	err := r.database.QueryRowContext(ctx, getTagByIDQuery, tagID).Scan(
		&tag.TagID,
		&tag.Name,
		&tag.Color,
		&tag.UserID,
		&tag.DeletedAt,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *TagRepository) Update(ctx context.Context, tag model.Tag) (*model.Tag, error) {
	var updated model.Tag
	err := r.database.QueryRowContext(
		ctx,
		updateTagQuery,
		tag.Name,
		tag.Color,
		tag.UpdatedAt,
		tag.TagID,
	).Scan(
		&updated.TagID,
		&updated.Name,
		&updated.Color,
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

func (r *TagRepository) Delete(ctx context.Context, tagID uuid.UUID, deletedAt time.Time) error {
	_, err := r.database.ExecContext(ctx, deleteTagQuery, deletedAt, tagID)
	if err != nil {
		return err
	}

	return nil
}

