package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	Preferences *UserPreferences `gorm:"foreignKey:UserID" json:"preferences,omitempty"`
	Threads     []Thread         `gorm:"foreignKey:UserID" json:"threads,omitempty"`
}
