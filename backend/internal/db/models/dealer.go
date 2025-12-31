package models

import (
	"time"

	"github.com/google/uuid"
)

type Dealer struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name             string    `gorm:"not null" json:"name"`
	Location         string    `gorm:"not null" json:"location"`
	Email            *string   `gorm:"type:varchar(255)" json:"email,omitempty"`
	Phone            *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Website          *string   `gorm:"type:varchar(255)" json:"website,omitempty"`
	Distance         float64   `gorm:"not null" json:"distance"`
	Contacted        bool      `gorm:"default:false" json:"contacted"`
	UserPreferenceID uuid.UUID `gorm:"type:uuid;index;not null" json:"userPreferenceId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	UserPreference *UserPreferences `gorm:"foreignKey:UserPreferenceID" json:"userPreference,omitempty"`
}
