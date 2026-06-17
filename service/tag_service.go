package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
)

type TagService struct {
	repository *repository.TagRepository
}

func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{
		repository: repo,
	}
}

func (srv *TagService) Create(ctx context.Context, request model.TagRequestDTO) (*model.Tag, error) {

	now := time.Now()

	entity := model.Tag{
		TagID:     uuid.New(),
		Name:      request.Name,
		Color:     request.Color,
		UserID:    request.UserID,
		DeletedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	saved_tag, err := srv.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return saved_tag, nil
}

func (srv *TagService) ListTagsFromUser(ctx context.Context, dto model.ListTagsFromUserDTO) ([]*model.Tag, error) {
	tags, err := srv.repository.ListTagsFromUser(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (srv *TagService) GetTagByID(ctx context.Context, tag_id uuid.UUID) (*model.Tag, error) {
	tag, err := srv.repository.GetTagByID(ctx, tag_id)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (srv *TagService) Update(ctx context.Context, tag_id uuid.UUID, request model.TagRequestDTO) (*model.Tag, error) {

	now := time.Now()

	entity := model.Tag{
		TagID:     tag_id,
		Name:      request.Name,
		Color:     request.Color,
		UpdatedAt: now,
	}

	updated_tag, err := srv.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	return updated_tag, nil
}

func (srv *TagService) Delete(ctx context.Context, tag_id uuid.UUID) error {

	now := time.Now()

	if err := srv.repository.Delete(ctx, tag_id, now); err != nil {
		return err
	}

	return nil
}