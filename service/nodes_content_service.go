package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
)

type NodesContentService struct {
	repository *repository.NodesContentRepository
}

func NewNodesContentService(repo *repository.NodesContentRepository) *NodesContentService {
	return &NodesContentService{
		repository: repo,
	}
}

func (srv *NodesContentService) Create(ctx context.Context, request model.NodesContentRequestDTO) (*model.NodesContent, error) {

	now := time.Now()

	entity := model.NodesContent{
		NodeID:     uuid.New(),
		ContentID:  request.ContentID,
		UserID:     request.UserID,
		NotebookID: request.NotebookID,
		DeletedAt:  nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	saved_node, err := srv.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return saved_node, nil
}

func (srv *NodesContentService) ListNodesContentsFromNotebook(ctx context.Context, dto model.ListNodesContentsFromNotebookDTO) ([]*model.NodesContent, error) {
	nodes, err := srv.repository.ListNodesContentsFromNotebook(ctx, dto.NotebookID)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (srv *NodesContentService) ListNodesContentsFromContent(ctx context.Context, dto model.ListNodesContentsFromContentDTO) ([]*model.NodesContent, error) {
	nodes, err := srv.repository.ListNodesContentsFromContent(ctx, dto.ContentID)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (srv *NodesContentService) GetNodesContentByID(ctx context.Context, node_id uuid.UUID) (*model.NodesContent, error) {
	node, err := srv.repository.GetNodesContentByID(ctx, node_id)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (srv *NodesContentService) Update(ctx context.Context, node_id uuid.UUID, request model.UpdateNodesContentDTO) (*model.NodesContent, error) {

	now := time.Now()

	entity := model.NodesContent{
		NodeID:    node_id,
		ContentID: request.ContentID,
		UpdatedAt: now,
	}

	updated_node, err := srv.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	return updated_node, nil
}

func (srv *NodesContentService) Delete(ctx context.Context, node_id uuid.UUID) error {

	now := time.Now()

	if err := srv.repository.Delete(ctx, node_id, now); err != nil {
		return err
	}

	return nil
}
