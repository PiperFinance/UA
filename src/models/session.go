package models

import "github.com/google/uuid"

type Session struct {
	UUID      *uuid.UUID `json:"UUID,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ExpiresAt int64
	UserRefer *uuid.UUID
	User      User `gorm:"foreignKey:UserRefer"`
}
