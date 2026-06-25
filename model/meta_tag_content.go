package model

import (
	"time"

	"github.com/google/uuid"
)

type MetaTagContent struct {
	ContentID  uuid.UUID  `json:"content_id"`
	TagID      uuid.UUID  `json:"tag_id"`
	NotebookID uuid.UUID  `json:"notebook_id"`
	UserID     uuid.UUID  `json:"user_id"`
	DeletedAt  *time.Time `json:"deleted_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type MetaTagContentRequestDTO struct {
	UserID     uuid.UUID
	NotebookID uuid.UUID
	ContentID  uuid.UUID `json:"content_id"`
	TagID      uuid.UUID `json:"tag_id"`
}

type UpdateMetaTagContentDTO struct {
	NotebookID uuid.UUID `json:"notebook_id"`
}

type ListMetaTagContentsFromContentDTO struct {
	ContentID uuid.UUID `json:"content_id"`
}

type ListMetaTagContentsFromTagDTO struct {
	TagID uuid.UUID `json:"tag_id"`
}

type ListMetaTagContentsFromNotebookDTO struct {
	NotebookID uuid.UUID `json:"notebook_id"`
}
