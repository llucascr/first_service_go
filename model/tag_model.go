package model

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	TagID     uuid.UUID  `json:"tag_id"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	UserID    uuid.UUID  `json:"user_id"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type TagRequestDTO struct {
	UserID uuid.UUID
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type ListTagsFromUserDTO struct {
	UserID uuid.UUID `json:"user_id"`
}
