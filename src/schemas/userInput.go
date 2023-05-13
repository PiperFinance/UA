package schemas

import (
	"github.com/PiperFinance/UA/src/models"
	"github.com/google/uuid"
)

type SignUpInput struct {
	Name      string         `json:"name"`
	Address   models.Address `json:"address" validate:"required"`
	Email     string         `json:"email"`
	Password  string         `json:"password" validate:"required,min=8"`
	SignedMsg string         `json:"signedMsg" validate:"hexadecimal,required"`
}

type SignInInput struct {
	Address  models.Address `json:"address"  validate:"required"`
	Password string         `json:"password"  validate:"required"`
	UserUUID uuid.UUID      `json:"uuid,omitempty" `
}
