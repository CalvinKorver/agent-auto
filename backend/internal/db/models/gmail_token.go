package models

import (
	"time"

	"github.com/google/uuid"
)

type GmailToken struct {
	UserID       uuid.UUID `gorm:"type:uuid;primary_key" json:"userId"`
	AccessToken  string    `gorm:"type:text;not null" json:"-"` // Encrypted, never expose in JSON
	RefreshToken string    `gorm:"type:text;not null" json:"-"` // Encrypted, never expose in JSON
	TokenType    string    `json:"tokenType,omitempty"`
	Expiry       time.Time `gorm:"not null" json:"expiry"`
	GmailEmail   string    `json:"gmailEmail,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
