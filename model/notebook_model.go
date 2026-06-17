package model

import (
	"time"

	"github.com/google/uuid"
)

type Notebook struct {
	NotebookID  uuid.UUID  `json:"notebook_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Icon        string     `json:"icon"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type NotebookRequestDTO struct {
	UserID      uuid.UUID
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type NotebookResponseDTO struct {
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type ListNotebooksFromUserDTO struct {
	UserID uuid.UUID `json:"user_id"`
}
