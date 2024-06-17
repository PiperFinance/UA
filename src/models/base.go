package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	UUID      *uuid.UUID     `json:"UUID,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UpdatedAt time.Time      `json:"createdAt" gorm:"autoUpdateTime:milli"`
	CreatedAt time.Time      `json:"updatedAt" gorm:"autoCreateTime:second"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (obj *BaseModel) ID() string {
	return obj.UUID.String()
}

// BeforeCreate will set a UUID rather than numeric ID.
//func (obj *BaseModel) BeforeCreate(scope *gorm.DB) error {
//	id, err := uuid.NewUUID()
//	if err == nil {
//		obj.UUID = &id
//	}
//	return err
//}
