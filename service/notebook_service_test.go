package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/llucascr/first_service_go/config"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/repository"
	"github.com/llucascr/first_service_go/service"
)


func createUser(db *sql.DB, id uuid.UUID, name, email, phone string) error {
	query := `INSERT INTO users (user_id, name, email, phone, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, id, name, email, phone, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func buildNotebookService() (*service.NoteBookService, uuid.UUID) {
	log.Println("creating a webserver")

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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
	return notebookService, userID
}

func TestNotebookService(t *testing.T) {
	notebookSrv, userID := buildNotebookService()
	if notebookSrv == nil {
		t.Error("The Service was not created")
	}

	// == Create ==
	t.Run("Create Notebook", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NotebookRequestDTO{
			UserID:      userID,
			Icon:        "folder",
			Name:        "Test Notebook",
			Image:       "google.com/image",
			Description: "Test Content",
		}

		output, err := notebookSrv.Create(ctx, input)
		require.NoError(t, err)

		t.Run("It should an id", func(t *testing.T) {
			require.NotNil(t, output.NotebookID)
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
	t.Run("List Notebooks from User", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NotebookRequestDTO{
			UserID:      userID,
			Icon:        "folder",
			Name:        "Test Notebook",
			Image:       "google.com/image",
			Description: "Test Content",
		}

		_, err := notebookSrv.Create(ctx, input)
		require.NoError(t, err)

		notebooks, err := notebookSrv.ListNotebooksFromUser(ctx, model.ListNotebooksFromUserDTO{
			UserID: userID,
		})
		require.NoError(t, err)

		t.Run("there should be 2 users", func(t *testing.T) {
			require.Len(t, notebooks, 2)
		})

	})

	t.Run("Get Notebook By ID", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NotebookRequestDTO{
			UserID:      userID,
			Icon:        "folder",
			Name:        "Test Notebook",
			Image:       "google.com/image",
			Description: "Test Content",
		}

		new_notebook, err := notebookSrv.Create(ctx, input)
		require.NoError(t, err)

		notebook_found, err := notebookSrv.GetNotebookByID(ctx,new_notebook.NotebookID)
		require.NoError(t, err)

		t.Run("the id should be equal", func(t *testing.T) {
			require.Equal(t, new_notebook.NotebookID, notebook_found.NotebookID)
		})

	})

	// == Update ==
	t.Run("Update Notebook", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NotebookRequestDTO{
			UserID:      userID,
			Icon:        "folder",
			Name:        "Test Notebook",
			Image:       "google.com/image",
			Description: "Test Content",
		}

		new_notebook, err := notebookSrv.Create(ctx, input)
		require.NoError(t, err)

		update_input := model.NotebookRequestDTO{
			Icon:        "book",
			Name:        "Updated Notebook",
			Image:       "google.com/updated",
			Description: "Updated Content",
		}

		updated_notebook, err := notebookSrv.Update(ctx, new_notebook.NotebookID, update_input)
		require.NoError(t, err)

		t.Run("the id should be the same", func(t *testing.T) {
			require.Equal(t, new_notebook.NotebookID, updated_notebook.NotebookID)
		})

		t.Run("the fields should be updated", func(t *testing.T) {
			require.Equal(t, "book", updated_notebook.Icon)
			require.Equal(t, "Updated Notebook", updated_notebook.Name)
			require.Equal(t, "google.com/updated", updated_notebook.Image)
			require.Equal(t, "Updated Content", updated_notebook.Description)
		})

		t.Run("it should update updated_at", func(t *testing.T) {
			require.NotEmpty(t, updated_notebook.UpdatedAt)
		})

	})

	// == Delete ==
	t.Run("Delete Notebook", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NotebookRequestDTO{
			UserID:      userID,
			Icon:        "folder",
			Name:        "Test Notebook",
			Image:       "google.com/image",
			Description: "Test Content",
		}

		new_notebook, err := notebookSrv.Create(ctx, input)
		require.NoError(t, err)

		err = notebookSrv.Delete(ctx, new_notebook.NotebookID)
		require.NoError(t, err)

		t.Run("it should not be found after delete", func(t *testing.T) {
			_, err := notebookSrv.GetNotebookByID(ctx, new_notebook.NotebookID)
			require.Error(t, err)
		})

	})
}
