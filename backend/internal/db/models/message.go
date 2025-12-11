package models

import (
	"time"

	"github.com/google/uuid"
)

type SenderType string

const (
	SenderTypeUser   SenderType = "user"
	SenderTypeAgent  SenderType = "agent"
	SenderTypeSeller SenderType = "seller"
)

type Message struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ThreadID  uuid.UUID  `gorm:"type:uuid;index;not null" json:"threadId"`
	Sender    SenderType `gorm:"type:varchar(20);not null" json:"sender"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	Timestamp time.Time  `gorm:"not null" json:"timestamp"`

	Thread *Thread `gorm:"foreignKey:ThreadID" json:"thread,omitempty"`
}
