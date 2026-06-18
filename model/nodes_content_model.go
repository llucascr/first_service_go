package model

import (
	"time"

	"github.com/google/uuid"
)

type NodesContent struct {
	NodeID     uuid.UUID  `json:"node_id"`
	ContentID  uuid.UUID  `json:"content_id"`
	UserID     uuid.UUID  `json:"user_id"`
	NotebookID uuid.UUID  `json:"notebook_id"`
	DeletedAt  *time.Time `json:"deleted_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type NodesContentRequestDTO struct {
	UserID     uuid.UUID
	NotebookID uuid.UUID
	ContentID  uuid.UUID `json:"content_id"`
}

type UpdateNodesContentDTO struct {
	ContentID uuid.UUID `json:"content_id"`
}

type ListNodesContentsFromNotebookDTO struct {
	NotebookID uuid.UUID `json:"notebook_id"`
}

type ListNodesContentsFromContentDTO struct {
	ContentID uuid.UUID `json:"content_id"`
}
