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
	ID                uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID            uuid.UUID   `gorm:"type:uuid;index;not null" json:"userId"`
	ThreadID          *uuid.UUID  `gorm:"type:uuid;index" json:"threadId,omitempty"`
	MessageTypeID     *uuid.UUID  `gorm:"type:uuid;index" json:"messageTypeId,omitempty"`
	Sender            SenderType `gorm:"type:varchar(20);not null" json:"sender"`
	Content           string      `gorm:"type:text;not null" json:"content"`
	Timestamp         time.Time   `gorm:"not null" json:"timestamp"`
	SenderEmail       string      `json:"senderEmail,omitempty"`
	SenderPhone       string      `json:"senderPhone,omitempty"`
	ExternalMessageID string      `gorm:"index" json:"externalMessageId,omitempty"`
	Subject           string      `json:"subject,omitempty"`
	Metadata          *string     `gorm:"type:jsonb" json:"metadata,omitempty"`
	SentViaEmail      bool        `gorm:"default:false" json:"sentViaEmail"`
	SentViaSMS        bool        `gorm:"default:false" json:"sentViaSMS"`
	DeletedAt         *time.Time  `gorm:"index" json:"deletedAt,omitempty"`

	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Thread      *Thread      `gorm:"foreignKey:ThreadID" json:"thread,omitempty"`
	MessageType *MessageType `gorm:"foreignKey:MessageTypeID" json:"messageType,omitempty"`
}
