package models

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MakeID    uuid.UUID `gorm:"type:uuid;index;not null" json:"makeId"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Make         *Make          `gorm:"foreignKey:MakeID" json:"make,omitempty"`
	VehicleTrims []VehicleTrim  `gorm:"foreignKey:ModelID" json:"vehicleTrims,omitempty"`
}

