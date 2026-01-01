package models

import (
	"github.com/google/uuid"
)

type MessageTypeEnum string

const (
	MessageTypeEmail MessageTypeEnum = "EMAIL"
	MessageTypePhone MessageTypeEnum = "PHONE"
)

type MessageType struct {
	ID   uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type MessageTypeEnum `gorm:"type:varchar(20);not null;unique" json:"type"`
}

