package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var ErrNotETHAddress = errors.New("address is not in a evm-compatible chain")

type Address struct {
	Hash     string `gorm:"primaryKey" validate:"required"`
	Chain    int64  // 1 -> ETH based , negative values for non - evm chains
	LastSync sql.NullTime
	Users    []*User `gorm:"many2many:user_addresses;"`
}

// ETHAddress error if address is not evm compatible ...
func (a Address) ETHAddress() (common.Address, error) {
	if common.IsHexAddress(a.Hash) {
		return common.HexToAddress(a.Hash), nil
	} else {
		return common.Address{}, ErrNotETHAddress
	}
}

// TODO - DO the same for other chains

func (a Address) VerifySignedMsg(OrgStr string, SignedMsg string) (bool, error) {
	signature, err := hexutil.Decode(SignedMsg)
	if err != nil {
		return false, fmt.Errorf("decode signature: %w", err)
	}
	return a.VerifySignedBytes([]byte(OrgStr), signature)
}

func (a Address) VerifySignedBytes(msg []byte, sig []byte) (bool, error) {
	msg = accounts.TextHash(msg)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	from, err := a.ETHAddress()
	if err != nil {
		// TODO - Implemnentaion is left for other chains ...
		return false, err
	}
	return from.Hex() == recoveredAddr.Hex(), nil
}
