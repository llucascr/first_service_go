package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
)

type MetaTagContentService struct {
	repository *repository.MetaTagContentRepository
}

func NewMetaTagContentService(repo *repository.MetaTagContentRepository) *MetaTagContentService {
	return &MetaTagContentService{
		repository: repo,
	}
}

func (srv *MetaTagContentService) Create(ctx context.Context, request model.MetaTagContentRequestDTO) (*model.MetaTagContent, error) {

	now := time.Now()

	entity := model.MetaTagContent{
		ContentID:  request.ContentID,
		TagID:      request.TagID,
		UserID:     request.UserID,
		NotebookID: request.NotebookID,
		DeletedAt:  nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	saved_meta_tag, err := srv.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return saved_meta_tag, nil
}

func (srv *MetaTagContentService) ListMetaTagContentsFromContent(ctx context.Context, dto model.ListMetaTagContentsFromContentDTO) ([]*model.MetaTagContent, error) {
	metaTags, err := srv.repository.ListMetaTagContentsFromContent(ctx, dto.ContentID)
	if err != nil {
		return nil, err
	}

	return metaTags, nil
}

func (srv *MetaTagContentService) ListMetaTagContentsFromTag(ctx context.Context, dto model.ListMetaTagContentsFromTagDTO) ([]*model.MetaTagContent, error) {
	metaTags, err := srv.repository.ListMetaTagContentsFromTag(ctx, dto.TagID)
	if err != nil {
		return nil, err
	}

	return metaTags, nil
}

func (srv *MetaTagContentService) ListMetaTagContentsFromNotebook(ctx context.Context, dto model.ListMetaTagContentsFromNotebookDTO) ([]*model.MetaTagContent, error) {
	metaTags, err := srv.repository.ListMetaTagContentsFromNotebook(ctx, dto.NotebookID)
	if err != nil {
		return nil, err
	}

	return metaTags, nil
}

func (srv *MetaTagContentService) GetMetaTagContent(ctx context.Context, content_id uuid.UUID, tag_id uuid.UUID) (*model.MetaTagContent, error) {
	metaTag, err := srv.repository.GetMetaTagContent(ctx, content_id, tag_id)
	if err != nil {
		return nil, err
	}

	return metaTag, nil
}

func (srv *MetaTagContentService) Update(ctx context.Context, content_id uuid.UUID, tag_id uuid.UUID, request model.UpdateMetaTagContentDTO) (*model.MetaTagContent, error) {

	now := time.Now()

	entity := model.MetaTagContent{
		ContentID:  content_id,
		TagID:      tag_id,
		NotebookID: request.NotebookID,
		UpdatedAt:  now,
	}

	updated_meta_tag, err := srv.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	return updated_meta_tag, nil
}

func (srv *MetaTagContentService) Delete(ctx context.Context, content_id uuid.UUID, tag_id uuid.UUID) error {

	now := time.Now()

	if err := srv.repository.Delete(ctx, content_id, tag_id, now); err != nil {
		return err
	}

	return nil
}
