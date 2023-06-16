package schemas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type (
	Id       uint32
	ChainId  int64
	TokenDet struct {
		ChainId     int64          `json:"chainId"`
		Address     common.Address `json:"address"`
		Name        string         `json:"name"`
		Symbol      string         `json:"symbol"`
		Decimals    int32          `json:"decimals"`
		Tags        []string       `json:"tags"`
		CoingeckoId string         `json:"coingeckoId"`
		LifiId      string         `json:"lifiId,omitempty"`
		ListedIn    []string       `json:"listedIn"`
		LogoURI     string         `json:"logoURI"`
		Verify      bool           `json:"verify"`
		Related     []Token        `json:"token,omitempty"`
	}

	Token struct {
		Detail              TokenDet  `json:"detail"`
		PriceUSD            float64   `json:"priceUSD"`
		Balance             big.Float `json:"-"`
		Value               big.Float `json:"-"`
		BalanceStr          string    `json:"balance"`
		BalanceNoDecimalStr string    `json:"balanceNoDecimal"`
		ValueStr            string    `json:"value"`
	}

	TokenId      Id
	TokenMapping map[TokenId]Token
)
