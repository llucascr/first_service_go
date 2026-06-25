package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
)

type MetaContentService struct {
	repository *repository.MetaContentRepository
}

func NewMetaContentService(repo *repository.MetaContentRepository) *MetaContentService {
	return &MetaContentService{
		repository: repo,
	}
}

func (srv *MetaContentService) Create(ctx context.Context, request model.MetaContentRequestDTO) (*model.MetaContent, error) {

	now := time.Now()

	entity := model.MetaContent{
		ContentID:  uuid.New(),
		NotebookID: request.NotebookID,
		UserID:     request.UserID,
		Icon:       request.Icon,
		Name:       request.Name,
		DeletedAt:  nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	saved_content, err := srv.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return saved_content, nil
}

func (srv *MetaContentService) ListMetaContentsFromNotebook(ctx context.Context, dto model.ListMetaContentsFromNotebookDTO) ([]*model.MetaContent, error) {
	contents, err := srv.repository.ListMetaContentsFromNotebook(ctx, dto.NotebookID)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (srv *MetaContentService) GetMetaContentByID(ctx context.Context, content_id uuid.UUID) (*model.MetaContent, error) {
	content, err := srv.repository.GetMetaContentByID(ctx, content_id)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (srv *MetaContentService) Update(ctx context.Context, content_id uuid.UUID, request model.UpdateMetaContentDTO) (*model.MetaContent, error) {

	now := time.Now()

	entity := model.MetaContent{
		ContentID: content_id,
		Icon:      request.Icon,
		Name:      request.Name,
		UpdatedAt: now,
	}

	updated_content, err := srv.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	return updated_content, nil
}

func (srv *MetaContentService) Delete(ctx context.Context, content_id uuid.UUID) error {

	now := time.Now()

	if err := srv.repository.Delete(ctx, content_id, now); err != nil {
		return err
	}

	return nil
}
