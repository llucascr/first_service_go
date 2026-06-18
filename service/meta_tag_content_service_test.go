package service_test

import (
	"context"
	"database/sql"
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

// createContentAndTag cria um meta_content e uma tag novos e devolve seus IDs,
// usado para isolar cenarios de Create da tabela de juncao.
func createContentAndTag(db *sql.DB, userID uuid.UUID, notebookID uuid.UUID) (uuid.UUID, uuid.UUID) {
	ctx := context.TODO()

	metaContentService := service.NewMetaContentService(repository.NewMetaContentRepository(db))
	content, err := metaContentService.Create(ctx, model.MetaContentRequestDTO{
		UserID:     userID,
		NotebookID: notebookID,
		Icon:       "file",
		Name:       "Fresh content",
	})
	if err != nil {
		panic(err)
	}

	tagService := service.NewTagService(repository.NewTagRepository(db))
	tag, err := tagService.Create(ctx, model.TagRequestDTO{
		UserID: userID,
		Name:   "Fresh tag",
		Color:  "ABCDEF",
	})
	if err != nil {
		panic(err)
	}

	return content.ContentID, tag.TagID
}

func buildMetaTagContentService() (*service.MetaTagContentService, *sql.DB, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID) {
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
	err = createUser(db, userID, "Test user", "metatagcontent@email.com", "1527298624")
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
		Description: "Notebook for meta tag content tests",
	})
	if err != nil {
		panic(err)
	}

	metaContentRepository := repository.NewMetaContentRepository(db)
	metaContentService := service.NewMetaContentService(metaContentRepository)
	content, err := metaContentService.Create(context.TODO(), model.MetaContentRequestDTO{
		UserID:     userID,
		NotebookID: notebook.NotebookID,
		Icon:       "file",
		Name:       "Content for meta tags",
	})
	if err != nil {
		panic(err)
	}

	tagRepository := repository.NewTagRepository(db)
	tagService := service.NewTagService(tagRepository)
	tag, err := tagService.Create(context.TODO(), model.TagRequestDTO{
		UserID: userID,
		Name:   "Tag for meta content",
		Color:  "FF0000",
	})
	if err != nil {
		panic(err)
	}

	metaTagContentRepository := repository.NewMetaTagContentRepository(db)
	metaTagContentService := service.NewMetaTagContentService(metaTagContentRepository)
	return metaTagContentService, db, userID, notebook.NotebookID, content.ContentID, tag.TagID
}

func TestMetaTagContentService(t *testing.T) {
	metaTagSrv, db, userID, notebookID, contentID, tagID := buildMetaTagContentService()
	if metaTagSrv == nil {
		t.Error("The Service was not created")
	}

	// == Create ==
	t.Run("Create Meta Tag Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.MetaTagContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  contentID,
			TagID:      tagID,
		}

		output, err := metaTagSrv.Create(ctx, input)
		require.NoError(t, err)

		t.Run("it should be linked to the content", func(t *testing.T) {
			require.Equal(t, contentID, output.ContentID)
		})

		t.Run("it should be linked to the tag", func(t *testing.T) {
			require.Equal(t, tagID, output.TagID)
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

	// == List from Content ==
	t.Run("List Meta Tag Contents from Content", func(t *testing.T) {
		ctx := context.TODO()

		metaTags, err := metaTagSrv.ListMetaTagContentsFromContent(ctx, model.ListMetaTagContentsFromContentDTO{
			ContentID: contentID,
		})
		require.NoError(t, err)

		t.Run("all should belong to the content", func(t *testing.T) {
			require.NotEmpty(t, metaTags)
			for _, metaTag := range metaTags {
				require.Equal(t, contentID, metaTag.ContentID)
			}
		})

	})

	// == List from Tag ==
	t.Run("List Meta Tag Contents from Tag", func(t *testing.T) {
		ctx := context.TODO()

		metaTags, err := metaTagSrv.ListMetaTagContentsFromTag(ctx, model.ListMetaTagContentsFromTagDTO{
			TagID: tagID,
		})
		require.NoError(t, err)

		t.Run("all should belong to the tag", func(t *testing.T) {
			require.NotEmpty(t, metaTags)
			for _, metaTag := range metaTags {
				require.Equal(t, tagID, metaTag.TagID)
			}
		})

	})

	// == List from Notebook ==
	t.Run("List Meta Tag Contents from Notebook", func(t *testing.T) {
		ctx := context.TODO()

		metaTags, err := metaTagSrv.ListMetaTagContentsFromNotebook(ctx, model.ListMetaTagContentsFromNotebookDTO{
			NotebookID: notebookID,
		})
		require.NoError(t, err)

		t.Run("all should belong to the notebook", func(t *testing.T) {
			require.NotEmpty(t, metaTags)
			for _, metaTag := range metaTags {
				require.Equal(t, notebookID, metaTag.NotebookID)
			}
		})

	})

	// == Get by composite key ==
	t.Run("Get Meta Tag Content", func(t *testing.T) {
		ctx := context.TODO()

		meta_tag_found, err := metaTagSrv.GetMetaTagContent(ctx, contentID, tagID)
		require.NoError(t, err)

		t.Run("the content_id should be equal", func(t *testing.T) {
			require.Equal(t, contentID, meta_tag_found.ContentID)
		})

		t.Run("the tag_id should be equal", func(t *testing.T) {
			require.Equal(t, tagID, meta_tag_found.TagID)
		})

	})

	// == Update ==
	t.Run("Update Meta Tag Content", func(t *testing.T) {
		ctx := context.TODO()

		// novo notebook para re-apontar a associacao
		notebookRepository := repository.NewNotebookRepository(db)
		notebookService := service.NewNotebookService(notebookRepository)
		other_notebook, err := notebookService.Create(ctx, model.NotebookRequestDTO{
			UserID:      userID,
			Icon:        "book",
			Name:        "Other Notebook",
			Image:       "",
			Description: "Other notebook for update test",
		})
		require.NoError(t, err)

		update_input := model.UpdateMetaTagContentDTO{
			NotebookID: other_notebook.NotebookID,
		}

		updated_meta_tag, err := metaTagSrv.Update(ctx, contentID, tagID, update_input)
		require.NoError(t, err)

		t.Run("the content_id should be the same", func(t *testing.T) {
			require.Equal(t, contentID, updated_meta_tag.ContentID)
		})

		t.Run("the tag_id should be the same", func(t *testing.T) {
			require.Equal(t, tagID, updated_meta_tag.TagID)
		})

		t.Run("the notebook_id should be updated", func(t *testing.T) {
			require.Equal(t, other_notebook.NotebookID, updated_meta_tag.NotebookID)
		})

		t.Run("it should update updated_at", func(t *testing.T) {
			require.NotEmpty(t, updated_meta_tag.UpdatedAt)
		})

	})

	// == Delete ==
	t.Run("Delete Meta Tag Content", func(t *testing.T) {
		ctx := context.TODO()

		err := metaTagSrv.Delete(ctx, contentID, tagID)
		require.NoError(t, err)

		t.Run("it should not be found after delete", func(t *testing.T) {
			_, err := metaTagSrv.GetMetaTagContent(ctx, contentID, tagID)
			require.Error(t, err)
		})

	})

	// == Re-create after soft delete (revive) ==
	t.Run("Re-create after soft delete should revive", func(t *testing.T) {
		ctx := context.TODO()

		freshContentID, freshTagID := createContentAndTag(db, userID, notebookID)

		input := model.MetaTagContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  freshContentID,
			TagID:      freshTagID,
		}

		_, err := metaTagSrv.Create(ctx, input)
		require.NoError(t, err)

		err = metaTagSrv.Delete(ctx, freshContentID, freshTagID)
		require.NoError(t, err)

		revived, err := metaTagSrv.Create(ctx, input)
		require.NoError(t, err)

		t.Run("it should not have deleted_at after revive", func(t *testing.T) {
			require.Empty(t, revived.DeletedAt)
		})

		t.Run("it should be found again after revive", func(t *testing.T) {
			found, err := metaTagSrv.GetMetaTagContent(ctx, freshContentID, freshTagID)
			require.NoError(t, err)
			require.Equal(t, freshContentID, found.ContentID)
			require.Equal(t, freshTagID, found.TagID)
		})

	})

	// == Create duplicate of an active association ==
	t.Run("Create duplicate active should error", func(t *testing.T) {
		ctx := context.TODO()

		freshContentID, freshTagID := createContentAndTag(db, userID, notebookID)

		input := model.MetaTagContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  freshContentID,
			TagID:      freshTagID,
		}

		_, err := metaTagSrv.Create(ctx, input)
		require.NoError(t, err)

		_, err = metaTagSrv.Create(ctx, input)
		require.Error(t, err)
	})
}
