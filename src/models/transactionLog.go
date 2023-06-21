package models

import "github.com/google/uuid"

type SwapRequest struct {
	BaseModel
	UserUUID    *uuid.UUID `json:"-" `
	User        User       `json:"user"`
	AddressHash string     `json:"-"`
	Address     Address    `json:"address" validate:"required"`
	ChainId     int64      `json:"chainId" validate:"required"`
	// NOTE - this following can be null

	FromTokenId          string `json:"fromTokenId" `
	ToTokenId            string `json:"toTokenId"`
	Amount               string `json:"amount"`
	SelectedService      string `json:"service"`
	SelectedSlippageRate string `json:"SlippageRate"`
	ClaimedAmountOut     string `json:"amountOut"`
	GivenRoute           string `json:"Route"`
	Data                 string `gorm:"type:text" json:"data"`
}
