package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPreferences struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"userId"`
	MakeID    uuid.UUID `gorm:"type:uuid;index;not null" json:"makeId"`
	ModelID   uuid.UUID `gorm:"type:uuid;index;not null" json:"modelId"`
	Year      int       `gorm:"not null" json:"year"`
	TrimID    *uuid.UUID `gorm:"type:uuid;index" json:"trimId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User  *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Make  *Make       `gorm:"foreignKey:MakeID" json:"make,omitempty"`
	Model *Model      `gorm:"foreignKey:ModelID" json:"model,omitempty"`
	Trim  *VehicleTrim `gorm:"foreignKey:TrimID" json:"trim,omitempty"`
}
