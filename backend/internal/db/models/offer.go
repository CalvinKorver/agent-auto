package models

import (
	"time"

	"github.com/google/uuid"
)

type TrackedOffer struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ThreadID  uuid.UUID  `gorm:"type:uuid;index;not null" json:"threadId"`
	MessageID *uuid.UUID `gorm:"type:uuid" json:"messageId,omitempty"`
	OfferText string     `gorm:"type:text;not null" json:"offerText"`
	TrackedAt time.Time  `gorm:"not null" json:"trackedAt"`

	Thread  *Thread  `gorm:"foreignKey:ThreadID" json:"thread,omitempty"`
	Message *Message `gorm:"foreignKey:MessageID" json:"message,omitempty"`
}
