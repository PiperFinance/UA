package models

import "database/sql"

// Device that user used to access website
type Device struct {
	BaseModel
	IP
	DeviceId sql.NullString `gorm:"uniqueIndex"`
}
