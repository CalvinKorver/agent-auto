package models

import (
	"time"

	"github.com/google/uuid"
)

type SellerType string

const (
	SellerTypePrivate    SellerType = "private"
	SellerTypeDealership SellerType = "dealership"
	SellerTypeOther      SellerType = "other"
)

type Thread struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;index;not null" json:"userId"`
	SellerName     string     `gorm:"not null" json:"sellerName"`
	SellerType     SellerType `gorm:"type:varchar(20);not null" json:"sellerType"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	LastMessageAt  *time.Time `json:"lastMessageAt,omitempty"`
	MessageCount   int        `gorm:"default:0" json:"messageCount"`
	DeletedAt      *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Messages      []Message      `gorm:"foreignKey:ThreadID" json:"messages,omitempty"`
	TrackedOffers []TrackedOffer `gorm:"foreignKey:ThreadID" json:"trackedOffers,omitempty"`
}
