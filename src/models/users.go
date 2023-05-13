package models

import (
	"database/sql"
	"github.com/google/uuid"
	"strings"
	"time"
)

type User struct {
	BaseModel
	Art         string       `json:"art,omitempty"`
	Name        string       `gorm:"type:varchar(100)" json:"name,omitempty"`
	Email       string       `gorm:"type:varchar(100);uniqueIndex;" json:"email,omitempty"`
	Age         uint8        `json:"age,omitempty"`
	Birthday    time.Time    `json:"birthday,omitempty"`
	ActivatedAt sql.NullTime `json:"activatedAt,omitempty"`
	Password    string       `gorm:"type:varchar(100);not null" json:"-"`
	Verified    *bool        `gorm:"not null;default:false"`
	Devices     []Device     `gorm:"many2many:user_devices;"`
	Addresses   []*Address   `gorm:"many2many:user_addresses;"`
	LastAccess  time.Time    `json:"LastAccess" `
}

func (u User) Add2Str() string {
	stringArray := make([]string, len(u.Addresses))
	for i, add := range u.Addresses {
		stringArray[i] = add.Hash
	}
	return strings.Join(stringArray, ",")
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		ID:        *user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
