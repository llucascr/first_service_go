package service_test

import (
	"context"
	"log"
	"testing"

	kivik "github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"

	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func build() *service.Service {
	client, err := kivik.New("couch", "http://admin:pass@localhost:5984/")
	if err != nil {
		log.Fatalf("Failed to create client: %s", err)
	}

	db := client.DB("notebooks")
	if err := db.Err(); err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)
	srv := service.NewService(repo)

	return srv
}

func TestService(t *testing.T) {
	srv := build()

	if srv == nil {
		t.Error("Servidor tem que existir")
	}

	t.Run("Service Create", func(t *testing.T) {
		ctx := context.TODO()
		input := model.NotebookRequestDTO{
			Title:   "Test Notebook",
			Content: "Test Content",
		}

		output, err := srv.Create(ctx, input)
		require.NoError(t, err)

		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, input.Content, output.Content)

		if output.ID == "" {
			t.Error("ID tem que existir no output")
		}
	})

	t.Run("Service Update", func(t *testing.T) {
		ctx := context.TODO()
		input := model.NotebookRequestDTO{
			Title:   "Test Updated Notebook",
			Content: "Test Updated Content",
		}

		output, err := srv.Update(ctx, input, "d8b31d1f-8bf6-4778-96a9-8f3ac16badca")
		require.NoError(t, err)

		assert.Equal(t, "d8b31d1f-8bf6-4778-96a9-8f3ac16badca", output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, input.Content, output.Content)
	})

	t.Run("Service Get", func(t *testing.T) {
		ctx := context.TODO()

		input_id := "d8b31d1f-8bf6-4778-96a9-8f3ac16badca"

		output, err := srv.Get(ctx, input_id)
		require.NoError(t, err)

		assert.Equal(t, "d8b31d1f-8bf6-4778-96a9-8f3ac16badca", output.ID)
		assert.Equal(t, output.Title, "Test Updated Notebook")
		assert.Equal(t, output.Content, "Test Updated Content")
	})

	t.Run("Service Delete", func(t *testing.T) {
		ctx := context.TODO()

		new_notebook := createNotebook(ctx, *srv)

		input_id := new_notebook.ID
		err := srv.Delete(ctx, input_id)
		require.NoError(t, err)

		_, err = srv.Get(ctx, input_id)
		require.Error(t, err)
	})

	t.Run("Service List", func(t *testing.T) {
		ctx := context.TODO()

		outputs, err := srv.List(ctx)
		require.NoError(t, err)

		assert.NotEmpty(t, len(outputs))
	})
}

func createNotebook(ctx context.Context, srv service.Service) *model.Notebook {
	input := model.NotebookRequestDTO{
		Title:   "Test Notebook",
		Content: "Test Content",
	}

	output, err := srv.Create(ctx, input)

	if err != nil {
		return nil
	}

	return output
}
