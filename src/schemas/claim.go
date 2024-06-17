package schemas

import (
	"strconv"
	"strings"

	"github.com/PiperFinance/UA/src/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaim struct {
	jwt.MapClaims
	submittedAddressCount int32     `json:"sba,omitempty"`
	addStr                string    `json:"adds,omitempty"`
	chainsStr             string    `json:"chns,omitempty"`
	UserUUID              uuid.UUID `json:"sub"`
	SessionUUID           uuid.UUID `json:"suid"`
	ExpiresAt             int64     `json:"exp"`
	IssuedAt              int64     `json:"iat"`
	NotBefore             int64     `json:"nbf"`
}

func (*TokenClaim) FromJwtMap(claims jwt.MapClaims) TokenClaim {
	// exp := (claims["suid"]).(string)
	return TokenClaim{
		UserUUID:    uuid.MustParse((claims["sub"]).(string)),
		SessionUUID: uuid.MustParse((claims["suid"]).(string)),
		// ExpiresAt: ,
	}
}

func (rt TokenClaim) Valid() error {
	return nil
}

func (rt TokenClaim) SubmittedAddressCount() int32 {
	return rt.submittedAddressCount
}

func (rt TokenClaim) AddStr() string {
	return rt.addStr
}

func (rt TokenClaim) ChainsStr() string {
	return rt.chainsStr
}

func (rt TokenClaim) SetAddresses(Addresses []*models.Address) {
	chainArray := make([]string, len(Addresses))
	addressArray := make([]string, len(Addresses))
	for i, add := range Addresses {
		addressArray[i] = add.Hash
		chainArray[i] = strconv.FormatInt(add.Chain, 10)
	}
	rt.addStr = strings.Join(addressArray, ",")
	rt.chainsStr = strings.Join(addressArray, ",")
	rt.submittedAddressCount = int32(len(Addresses))
}

func (rt TokenClaim) GetAddresses() ([]models.Address, error) {
	// chainArray := make([]string, rt.submittedAddressCount)
	// addressArray := make([]string, rt.submittedAddressCount)
	resArray := make([]models.Address, rt.submittedAddressCount)
	for i, add := range strings.Split(rt.addStr, ",") {
		resArray[i].Hash = add
	}
	for i, chain := range strings.Split(rt.chainsStr, ",") {
		x, err := strconv.ParseInt(chain, 10, 64)
		resArray[i].Chain = x
		if err != nil {
			return nil, err
		}
	}
	return resArray, nil
}
