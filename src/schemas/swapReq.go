package schemas

type SwapReq struct {
	AddressHash          string `json:"addressHash" validate:"required"`
	ChainId              int64  `json:"chainId" validate:"required"`
	FromTokenId          string `json:"fromTokenId" validate:"required"`
	ToTokenId            string `json:"toTokenId" validate:"required"`
	Amount               string `json:"amount"`
	SelectedService      string `json:"service"`
	SelectedSlippageRate string `json:"SlippageRate"`
	ClaimedAmountOut     string `json:"amountOut"`
	GivenRoute           string `json:"Route"`
	Data                 string `json:"data"`
}
