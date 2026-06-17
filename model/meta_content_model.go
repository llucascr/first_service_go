package model

import (
	"time"

	"github.com/google/uuid"
)

type MetaContent struct {
	ContentID  uuid.UUID  `json:"content_id"`
	NotebookID uuid.UUID  `json:"notebook_id"`
	UserID     uuid.UUID  `json:"user_id"`
	Icon       string     `json:"icon"`
	Name       string     `json:"name"`
	DeletedAt  *time.Time `json:"deleted_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type MetaContentRequestDTO struct {
	UserID     uuid.UUID
	NotebookID uuid.UUID
	Icon       string `json:"icon"`
	Name       string `json:"name"`
}

type UpdateMetaContentDTO struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

type ListMetaContentsFromNotebookDTO struct {
	NotebookID uuid.UUID `json:"notebook_id"`
}