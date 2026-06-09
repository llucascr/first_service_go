package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/model"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

// Create
func (srv *Service) Create(ctx context.Context, request model.NotebookRequestDTO) (*model.Notebook, error) {

	new_notebook := &model.Notebook{
		ID: uuid.NewString(),
		Title: request.Title,
		Content: request.Content,
	}

	err := srv.repository.Create(ctx, new_notebook)
	if err != nil {
		return nil, err
	}

	return new_notebook, nil
}

// Update
func (srv *Service) Update(ctx context.Context, request model.NotebookRequestDTO, notebook_id string) (*model.Notebook, error) {

	notebook_updated, err := srv.repository.Update(ctx, notebook_id, request)
	if err != nil {
		return nil, err
	}
	
	return notebook_updated, nil
}

// Get
func (srv *Service) Get(ctx context.Context, notebook_id string) (*model.Notebook, error)  {

	notebook, err := srv.repository.Get(ctx, notebook_id)
	if err != nil {
		return nil, err
	}

	return notebook, nil
}

// Delete
func (srv *Service) Delete(ctx context.Context, notebook_id string) error {
	
	err := srv.repository.Delete(ctx, notebook_id)
	if err != nil {
		return err
	}

	return nil
}

// List
func (srv *Service) List(ctx context.Context) ([]*model.Notebook, error){

	notebooks, err := srv.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return notebooks, nil
}
