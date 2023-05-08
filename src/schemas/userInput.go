package schemas

import "github.com/PiperFinance/UA/src/models"

type SignUpInput struct {
	Name     string         `json:"name"`
	Address  models.Address `json:"address" `
	Email    string         `json:"email"`
	Password string         `json:"password" validate:"required,min=8"`
}

type SignInInput struct {
	Address  models.Address `json:"address"  validate:"required"`
	Password string         `json:"password"  validate:"required"`
}
