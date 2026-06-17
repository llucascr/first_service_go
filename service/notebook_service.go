package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
)

type NoteBookService struct {
	repository *repository.NotebookRepository
}

func NewNotebookService(repo *repository.NotebookRepository) *NoteBookService {
	return &NoteBookService{
		repository: repo,
	}
}

func (srv *NoteBookService) Create(ctx context.Context, request model.NotebookRequestDTO) (*model.Notebook, error) {

	now := time.Now()

	entity := model.Notebook{
		NotebookID: uuid.New(),
		UserID: request.UserID,
		Icon: request.Icon,
		Name: request.Name,
		Image: request.Image,
		Description: request.Description,
		DeletedAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	saved_notebook, err := srv.repository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return saved_notebook, nil
}

func (srv *NoteBookService) ListNotebooksFromUser(ctx context.Context, dto model.ListNotebooksFromUserDTO) ([]*model.Notebook, error) {
	notebooks, err := srv.repository.ListNotebooksFromUser(ctx, dto.UserID)
	if err != nil {
		return nil, err
	}

	return notebooks, nil
}

func (srv *NoteBookService) GetNotebookByID(ctx context.Context, notebook_id uuid.UUID) (*model.Notebook, error) {
	notebook, err := srv.repository.GetNotebookByID(ctx, notebook_id)
	if err != nil {
		return nil, err
	}

	return notebook, nil
}

func (srv *NoteBookService) Update(ctx context.Context, notebook_id uuid.UUID, request model.NotebookRequestDTO) (*model.Notebook, error) {

	now := time.Now()

	entity := model.Notebook{
		NotebookID: notebook_id,
		Icon: request.Icon,
		Name: request.Name,
		Image: request.Image,
		Description: request.Description,
		UpdatedAt: now,
	}

	updated_notebook, err := srv.repository.Update(ctx, entity)
	if err != nil {
		return nil, err
	}

	return updated_notebook, nil
}

func (srv *NoteBookService) Delete(ctx context.Context, notebook_id uuid.UUID) error {

	now := time.Now()

	if err := srv.repository.Delete(ctx, notebook_id, now); err != nil {
		return err
	}

	return nil
}
