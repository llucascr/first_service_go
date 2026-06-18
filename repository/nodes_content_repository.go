package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/google/uuid"
	"github.com/llucascr/first_service_go/model"
)

type NodesContentRepository struct {
	database *sql.DB
}

//go:embed queries/nodes_content_create.sql
var createNodesContentQuery string

//go:embed queries/nodes_content_list_from_notebook.sql
var listNodesContentFromNotebookQuery string

//go:embed queries/nodes_content_list_from_content.sql
var listNodesContentFromContentQuery string

//go:embed queries/nodes_content_get_by_id.sql
var getNodesContentByIDQuery string

//go:embed queries/nodes_content_update.sql
var updateNodesContentQuery string

//go:embed queries/nodes_content_delete.sql
var deleteNodesContentQuery string

func NewNodesContentRepository(db *sql.DB) *NodesContentRepository {
	return &NodesContentRepository{
		database: db,
	}
}

func (r *NodesContentRepository) Create(ctx context.Context, node model.NodesContent) (*model.NodesContent, error) {
	_, err := r.database.ExecContext(
		ctx,
		createNodesContentQuery,
		node.NodeID,
		node.ContentID,
		node.UserID,
		node.NotebookID,
		node.DeletedAt,
		node.CreatedAt,
		node.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (r *NodesContentRepository) ListNodesContentsFromNotebook(ctx context.Context, notebookID uuid.UUID) ([]*model.NodesContent, error) {
	rows, err := r.database.QueryContext(ctx, listNodesContentFromNotebookQuery, notebookID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nodes []*model.NodesContent

	for rows.Next() {
		var node model.NodesContent
		if err := rows.Scan(
			&node.NodeID,
			&node.ContentID,
			&node.UserID,
			&node.NotebookID,
			&node.CreatedAt,
			&node.UpdatedAt,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *NodesContentRepository) ListNodesContentsFromContent(ctx context.Context, contentID uuid.UUID) ([]*model.NodesContent, error) {
	rows, err := r.database.QueryContext(ctx, listNodesContentFromContentQuery, contentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nodes []*model.NodesContent

	for rows.Next() {
		var node model.NodesContent
		if err := rows.Scan(
			&node.NodeID,
			&node.ContentID,
			&node.UserID,
			&node.NotebookID,
			&node.CreatedAt,
			&node.UpdatedAt,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *NodesContentRepository) GetNodesContentByID(ctx context.Context, nodeID uuid.UUID) (*model.NodesContent, error) {
	var node model.NodesContent
	err := r.database.QueryRowContext(ctx, getNodesContentByIDQuery, nodeID).Scan(
		&node.NodeID,
		&node.ContentID,
		&node.UserID,
		&node.NotebookID,
		&node.DeletedAt,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (r *NodesContentRepository) Update(ctx context.Context, node model.NodesContent) (*model.NodesContent, error) {
	var updated model.NodesContent
	err := r.database.QueryRowContext(
		ctx,
		updateNodesContentQuery,
		node.ContentID,
		node.UpdatedAt,
		node.NodeID,
	).Scan(
		&updated.NodeID,
		&updated.ContentID,
		&updated.UserID,
		&updated.NotebookID,
		&updated.DeletedAt,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *NodesContentRepository) Delete(ctx context.Context, nodeID uuid.UUID, deletedAt time.Time) error {
	_, err := r.database.ExecContext(ctx, deleteNodesContentQuery, deletedAt, nodeID)
	if err != nil {
		return err
	}

	return nil
}
