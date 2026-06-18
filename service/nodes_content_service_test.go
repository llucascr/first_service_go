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

func buildNodesContentService() (*service.NodesContentService, *service.MetaContentService, uuid.UUID, uuid.UUID, uuid.UUID) {
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
	err = createUser(db, userID, "Test user", "nodescontent@email.com", "1527298624")
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
		Description: "Notebook for nodes content tests",
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
		Name:       "Content for nodes",
	})
	if err != nil {
		panic(err)
	}

	nodesContentRepository := repository.NewNodesContentRepository(db)
	nodesContentService := service.NewNodesContentService(nodesContentRepository)
	return nodesContentService, metaContentService, userID, notebook.NotebookID, content.ContentID
}

func TestNodesContentService(t *testing.T) {
	nodeSrv, metaContentSrv, userID, notebookID, contentID := buildNodesContentService()
	if nodeSrv == nil {
		t.Error("The Service was not created")
	}

	// == Create ==
	t.Run("Create Nodes Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NodesContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  contentID,
		}

		output, err := nodeSrv.Create(ctx, input)
		require.NoError(t, err)

		t.Run("it should have a node id", func(t *testing.T) {
			require.NotNil(t, output.NodeID)
		})

		t.Run("it should belong to the notebook", func(t *testing.T) {
			require.Equal(t, notebookID, output.NotebookID)
		})

		t.Run("it should be linked to the content", func(t *testing.T) {
			require.Equal(t, contentID, output.ContentID)
		})

		t.Run("it should create_at and update_at", func(t *testing.T) {
			require.NotEmpty(t, output.CreatedAt)
			require.NotEmpty(t, output.UpdatedAt)
		})

		t.Run("it should not have deleted_at", func(t *testing.T) {
			require.Empty(t, output.DeletedAt)
		})

	})

	// == List from Notebook ==
	t.Run("List Nodes Contents from Notebook", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NodesContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  contentID,
		}

		_, err := nodeSrv.Create(ctx, input)
		require.NoError(t, err)

		nodes, err := nodeSrv.ListNodesContentsFromNotebook(ctx, model.ListNodesContentsFromNotebookDTO{
			NotebookID: notebookID,
		})
		require.NoError(t, err)

		t.Run("there should be 2 nodes", func(t *testing.T) {
			require.Len(t, nodes, 2)
		})

	})

	// == List from Content ==
	t.Run("List Nodes Contents from Content", func(t *testing.T) {
		ctx := context.TODO()

		nodes, err := nodeSrv.ListNodesContentsFromContent(ctx, model.ListNodesContentsFromContentDTO{
			ContentID: contentID,
		})
		require.NoError(t, err)

		t.Run("all nodes should belong to the content", func(t *testing.T) {
			require.NotEmpty(t, nodes)
			for _, node := range nodes {
				require.Equal(t, contentID, node.ContentID)
			}
		})

	})

	// == Get By ID ==
	t.Run("Get Nodes Content By ID", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NodesContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  contentID,
		}

		new_node, err := nodeSrv.Create(ctx, input)
		require.NoError(t, err)

		node_found, err := nodeSrv.GetNodesContentByID(ctx, new_node.NodeID)
		require.NoError(t, err)

		t.Run("the id should be equal", func(t *testing.T) {
			require.Equal(t, new_node.NodeID, node_found.NodeID)
		})

	})

	// == Update ==
	t.Run("Update Nodes Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NodesContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  contentID,
		}

		new_node, err := nodeSrv.Create(ctx, input)
		require.NoError(t, err)

		// novo meta content para re-apontar o vínculo
		other_content, err := metaContentSrv.Create(ctx, model.MetaContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			Icon:       "star",
			Name:       "Other Content",
		})
		require.NoError(t, err)

		update_input := model.UpdateNodesContentDTO{
			ContentID: other_content.ContentID,
		}

		updated_node, err := nodeSrv.Update(ctx, new_node.NodeID, update_input)
		require.NoError(t, err)

		t.Run("the node id should be the same", func(t *testing.T) {
			require.Equal(t, new_node.NodeID, updated_node.NodeID)
		})

		t.Run("the content_id should be updated", func(t *testing.T) {
			require.Equal(t, other_content.ContentID, updated_node.ContentID)
		})

		t.Run("it should update updated_at", func(t *testing.T) {
			require.NotEmpty(t, updated_node.UpdatedAt)
		})

	})

	// == Delete ==
	t.Run("Delete Nodes Content", func(t *testing.T) {
		ctx := context.TODO()

		input := model.NodesContentRequestDTO{
			UserID:     userID,
			NotebookID: notebookID,
			ContentID:  contentID,
		}

		new_node, err := nodeSrv.Create(ctx, input)
		require.NoError(t, err)

		err = nodeSrv.Delete(ctx, new_node.NodeID)
		require.NoError(t, err)

		t.Run("it should not be found after delete", func(t *testing.T) {
			_, err := nodeSrv.GetNodesContentByID(ctx, new_node.NodeID)
			require.Error(t, err)
		})

	})
}
