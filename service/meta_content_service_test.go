package service_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/llucascr/first_service_go/config"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)

func buildMetaContentService() (*service.MetaContentService, uuid.UUID, uuid.UUID) {
	log.Println("creating a webserver")

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := config.NewPostgresDB(connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to PostgreSQL successfully established!")

	userID := uuid.New()
	err = createUser(db, userID, "Test user", "test@email.com", "1527298624")
	if err != nil {
		panic(err)
	}

	notebookRepository := repository.NewNotebookRepository(db)
	notebookService := service.NewNotebookService(notebookRepository)
	notebook, err := notebookService.Create(context.TODO(), model.NotebookRequestDTO{
		UserID:      userID,
		Icon:        "book",
		Name:        "Test Notebook",
		Image:       "",
		Description: "Notebook for meta content tests",
	})
	if err != nil {
		panic(err)
	}

	metaContentRepository := repository.NewMetaContentRepository(db)
	metaContentService := service.NewMetaContentService(metaContentRepository)
	return metaContentService, userID, notebook.NotebookID
}

func TestMetaContentService(t *testing.T) {
	contentSrv, userID, notebookID := buildMetaContentService()
	if contentSrv == nil {
		t.Error("The Service was not created")
	}

	// == Create ==
	t.Run("Create Meta Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.MetaContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			Icon:       "file",
			Name:       "Test Content",
		}

		output, err := contentSrv.Create(ctx, input)
		require.NoError(t, err)

		t.Run("it should have an id", func(t *testing.T) {
			require.NotNil(t, output.ContentID)
		})

		t.Run("it should belong to the notebook", func(t *testing.T) {
			require.Equal(t, notebookID, output.NotebookID)
		})

		t.Run("it should create_at and update_at", func(t *testing.T) {
			require.NotEmpty(t, output.CreatedAt)
			require.NotEmpty(t, output.UpdatedAt)
		})

		t.Run("it should not have deleted_at", func(t *testing.T) {
			require.Empty(t, output.DeletedAt)
		})

	})

	// == List ==
	t.Run("List Meta Contents from Notebook", func(t *testing.T) {
		ctx := context.TODO()

		input := model.MetaContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			Icon:       "file",
			Name:       "Second Content",
		}

		_, err := contentSrv.Create(ctx, input)
		require.NoError(t, err)

		contents, err := contentSrv.ListMetaContentsFromNotebook(ctx, model.ListMetaContentsFromNotebookDTO{
			NotebookID: notebookID,
		})
		require.NoError(t, err)

		t.Run("there should be 2 contents", func(t *testing.T) {
			require.Len(t, contents, 2)
		})

	})

	// == Get By ID ==
	t.Run("Get Meta Content By ID", func(t *testing.T) {
		ctx := context.TODO()

		input := model.MetaContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			Icon:       "file",
			Name:       "Lookup Content",
		}

		new_content, err := contentSrv.Create(ctx, input)
		require.NoError(t, err)

		content_found, err := contentSrv.GetMetaContentByID(ctx, new_content.ContentID)
		require.NoError(t, err)

		t.Run("the id should be equal", func(t *testing.T) {
			require.Equal(t, new_content.ContentID, content_found.ContentID)
		})

	})

	// == Update ==
	t.Run("Update Meta Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.MetaContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			Icon:       "file",
			Name:       "Original Content",
		}

		new_content, err := contentSrv.Create(ctx, input)
		require.NoError(t, err)

		update_input := model.UpdateMetaContentDTO{
			Icon: "star",
			Name: "Updated Content",
		}

		updated_content, err := contentSrv.Update(ctx, new_content.ContentID, update_input)
		require.NoError(t, err)

		t.Run("the id should be the same", func(t *testing.T) {
			require.Equal(t, new_content.ContentID, updated_content.ContentID)
		})

		t.Run("the fields should be updated", func(t *testing.T) {
			require.Equal(t, "star", updated_content.Icon)
			require.Equal(t, "Updated Content", updated_content.Name)
		})

		t.Run("it should update updated_at", func(t *testing.T) {
			require.NotEmpty(t, updated_content.UpdatedAt)
		})

	})

	// == Delete ==
	t.Run("Delete Meta Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.MetaContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			Icon:       "file",
			Name:       "To Delete",
		}

		new_content, err := contentSrv.Create(ctx, input)
		require.NoError(t, err)

		err = contentSrv.Delete(ctx, new_content.ContentID)
		require.NoError(t, err)

		t.Run("it should not be found after delete", func(t *testing.T) {
			_, err := contentSrv.GetMetaContentByID(ctx, new_content.ContentID)
			require.Error(t, err)
		})

	})
}