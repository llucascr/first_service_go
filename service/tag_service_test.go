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

func buildTag() (*service.TagService, uuid.UUID) {
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

	tagRepository := repository.NewTagRepository(db)
	tagService := service.NewTagService(tagRepository)
	return tagService, userID
}

func TestTagService(t *testing.T) {
	tagSrv, userID := buildTag()
	if tagSrv == nil {
		t.Error("The Service was not created")
	}

	// == Create ==
	t.Run("Create Tag", func(t *testing.T) {
		ctx := context.TODO()

		input := model.TagRequestDTO{
			UserID: userID,
			Name:   "Test Tag",
			Color:  "FF0000",
		}

		output, err := tagSrv.Create(ctx, input)
		require.NoError(t, err)

		t.Run("It should an id", func(t *testing.T) {
			require.NotNil(t, output.TagID)
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
	t.Run("List Tags from User", func(t *testing.T) {
		ctx := context.TODO()

		input := model.TagRequestDTO{
			UserID: userID,
			Name:   "Test Tag",
			Color:  "00FF00",
		}

		_, err := tagSrv.Create(ctx, input)
		require.NoError(t, err)

		tags, err := tagSrv.ListTagsFromUser(ctx, model.ListTagsFromUserDTO{
			UserID: userID,
		})
		require.NoError(t, err)

		t.Run("there should be 2 tags", func(t *testing.T) {
			require.Len(t, tags, 2)
		})

	})

	t.Run("Get Tag By ID", func(t *testing.T) {
		ctx := context.TODO()

		input := model.TagRequestDTO{
			UserID: userID,
			Name:   "Test Tag",
			Color:  "0000FF",
		}

		new_tag, err := tagSrv.Create(ctx, input)
		require.NoError(t, err)

		tag_found, err := tagSrv.GetTagByID(ctx, new_tag.TagID)
		require.NoError(t, err)

		t.Run("the id should be equal", func(t *testing.T) {
			require.Equal(t, new_tag.TagID, tag_found.TagID)
		})

	})

	// == Update ==
	t.Run("Update Tag", func(t *testing.T) {
		ctx := context.TODO()

		input := model.TagRequestDTO{
			UserID: userID,
			Name:   "Test Tag",
			Color:  "FF0000",
		}

		new_tag, err := tagSrv.Create(ctx, input)
		require.NoError(t, err)

		update_input := model.TagRequestDTO{
			Name:  "Updated Tag",
			Color: "123456",
		}

		updated_tag, err := tagSrv.Update(ctx, new_tag.TagID, update_input)
		require.NoError(t, err)

		t.Run("the id should be the same", func(t *testing.T) {
			require.Equal(t, new_tag.TagID, updated_tag.TagID)
		})

		t.Run("the fields should be updated", func(t *testing.T) {
			require.Equal(t, "Updated Tag", updated_tag.Name)
			require.Equal(t, "123456", updated_tag.Color)
		})

		t.Run("it should update updated_at", func(t *testing.T) {
			require.NotEmpty(t, updated_tag.UpdatedAt)
		})

	})

	// == Delete ==
	t.Run("Delete Tag", func(t *testing.T) {
		ctx := context.TODO()

		input := model.TagRequestDTO{
			UserID: userID,
			Name:   "Test Tag",
			Color:  "FF0000",
		}

		new_tag, err := tagSrv.Create(ctx, input)
		require.NoError(t, err)

		err = tagSrv.Delete(ctx, new_tag.TagID)
		require.NoError(t, err)

		t.Run("it should not be found after delete", func(t *testing.T) {
			_, err := tagSrv.GetTagByID(ctx, new_tag.TagID)
			require.Error(t, err)
		})

	})
}
