package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPreferences struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	Make      string    `gorm:"not null" json:"make"`
	Model     string    `gorm:"not null" json:"model"`
	Year      int       `gorm:"not null" json:"year"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
