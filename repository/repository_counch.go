package repository

import (
	"context"
	"errors"

	kivik "github.com/go-kivik/kivik/v4"

	"github.com/llucascr/first_service_go/model"
)

type NotebookCounch struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Rev     string `json:"_rev,omitempty"`
}

type Repository struct {
	counchDB *kivik.DB
}

func NewRepository(c *kivik.DB) *Repository {
	return &Repository{
		counchDB: c,
	}
}

func (r *Repository) Create(ctx context.Context, notebook *model.Notebook) error {
	notebook_counch := NotebookCounch{
		ID:      notebook.ID,
		Title:   notebook.Title,
		Content: notebook.Content,
		Rev:     "",
	}

	_, err := r.counchDB.Put(ctx, notebook_counch.ID, notebook_counch)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, notebook_id string, notebook model.NotebookRequestDTO) (*model.Notebook, error) {

	row := r.counchDB.Get(ctx, notebook_id)

	var current map[string]interface{}
	if err := row.ScanDoc(&current); err != nil {
		return nil, err
	}

	current["title"] = notebook.Title
	current["content"] = notebook.Content

	_, err := r.counchDB.Put(ctx, notebook_id, current)
	if err != nil {
		return nil, err
	}

	return &model.Notebook{
		ID: notebook_id,
		Title: notebook.Title,
		Content: notebook.Content,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, notebook_id string) error {
	var nb NotebookCounch
	row := r.counchDB.Get(ctx, notebook_id)
	if err := row.ScanDoc(&nb); err != nil {
		return err
	}

	rev := nb.Rev
	if rev == "" {
		return errors.New("Revision not found")
	}

	_, err := r.counchDB.Delete(ctx, notebook_id, rev)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Notebook, error) {
	var notebook model.Notebook
	err := r.counchDB.Get(ctx, id).ScanDoc(&notebook)

	if err != nil {
		return nil, err
	}

	return &notebook, nil
}

func (r *Repository) List(ctx context.Context) ([]*model.Notebook, error) {

	rows := r.counchDB.AllDocs(
		ctx,
		kivik.Param("include_docs", true),
	)

	var notebooks []*model.Notebook

	for rows.Next() {
		notebook := &model.Notebook{}

		if err := rows.ScanDoc(notebook); err != nil {
			return nil, err
		}

		notebooks = append(notebooks, notebook)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notebooks, nil
}
