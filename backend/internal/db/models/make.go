package models

import (
	"time"

	"github.com/google/uuid"
)

type Make struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Models []Model `gorm:"foreignKey:MakeID" json:"models,omitempty"`
}

