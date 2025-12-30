package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type VehicleTrim struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ModelID   uuid.UUID `gorm:"type:uuid;index;not null" json:"modelId"`
	Year      int       `gorm:"index;not null" json:"year"`
	TrimName  string    `gorm:"index;not null" json:"trimName"`
	TrimDescription sql.NullString `gorm:"type:varchar(155)" json:"trimDescription,omitempty"`
	BaseMSRP  sql.NullFloat64 `gorm:"type:decimal(10,2)" json:"baseMsrp,omitempty"`
	ImageURLs sql.NullString `gorm:"type:text" json:"imageUrls,omitempty"`
	SourceID  sql.NullString `gorm:"type:varchar(9)" json:"sourceId,omitempty"`
	
	// Spec fields - all nullable
	BodyType            sql.NullString `gorm:"type:varchar(20)" json:"bodyType,omitempty"`
	LengthIn            sql.NullFloat64 `gorm:"type:decimal(5,1)" json:"lengthIn,omitempty"`
	WidthIn             sql.NullFloat64 `gorm:"type:decimal(4,1)" json:"widthIn,omitempty"`
	HeightIn            sql.NullFloat64 `gorm:"type:decimal(4,1)" json:"heightIn,omitempty"`
	WheelbaseIn         sql.NullFloat64 `gorm:"type:decimal(5,1)" json:"wheelbaseIn,omitempty"`
	CurbWeightLbs       sql.NullInt32 `json:"curbWeightLbs,omitempty"`
	Cylinders           sql.NullString `gorm:"type:varchar(8)" json:"cylinders,omitempty"`
	EngineSizeL         sql.NullFloat64 `gorm:"type:decimal(3,1)" json:"engineSizeL,omitempty"`
	HorsepowerHP        sql.NullInt32 `json:"horsepowerHp,omitempty"`
	HorsepowerRPM       sql.NullInt32 `json:"horsepowerRpm,omitempty"`
	TorqueFtLbs         sql.NullInt32 `json:"torqueFtLbs,omitempty"`
	TorqueRPM           sql.NullInt32 `json:"torqueRpm,omitempty"`
	DriveType           sql.NullString `gorm:"type:varchar(17)" json:"driveType,omitempty"`
	Transmission        sql.NullString `gorm:"type:varchar(47)" json:"transmission,omitempty"`
	EngineType          sql.NullString `gorm:"type:varchar(20)" json:"engineType,omitempty"`
	FuelType            sql.NullString `gorm:"type:varchar(44)" json:"fuelType,omitempty"`
	FuelTankCapacityGal sql.NullFloat64 `gorm:"type:decimal(4,1)" json:"fuelTankCapacityGal,omitempty"`
	EPACombinedMPG      sql.NullInt32 `json:"epaCombinedMpg,omitempty"`
	EPACityHighwayMPG   sql.NullString `gorm:"type:varchar(10)" json:"epaCityHighwayMpg,omitempty"`
	RangeMiles           sql.NullString `gorm:"type:varchar(17)" json:"rangeMiles,omitempty"`
	CountryOfOrigin      sql.NullString `gorm:"type:varchar(21)" json:"countryOfOrigin,omitempty"`
	CarClassification    sql.NullString `gorm:"type:varchar(24)" json:"carClassification,omitempty"`
	PlatformCode         sql.NullString `gorm:"type:varchar(48)" json:"platformCode,omitempty"`
	
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Model *Model `gorm:"foreignKey:ModelID" json:"model,omitempty"`
}

