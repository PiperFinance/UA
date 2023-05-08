package models

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var NotETHAddress = errors.New("address is not in a evm-compatible chain")

type Address struct {
	Hash  string `gorm:"primaryKey" validate:"required"`
	Chain int64
}

// ETHAddress error if address is not evm compatible ...
func (a Address) ETHAddress() (common.Address, error) {
	if common.IsHexAddress(a.Hash) {
		return common.HexToAddress(a.Hash), nil
	} else {
		return common.Address{}, NotETHAddress
	}
}

// TODO - DO the same for other chains

func (a Address) VerifySignedMsg(OrgStr string, SignedMsg string) (bool, error) {
	return a.VerifySignedBytes([]byte(OrgStr), []byte(SignedMsg))
}
func (a Address) VerifySignedBytes(OrgBytes []byte, SignedBytes []byte) (bool, error) {
	sigPublicKey, err := crypto.Ecrecover(OrgBytes, SignedBytes)
	if err != nil {
		return false, err
	}
	add, err := a.ETHAddress()
	if err != nil {
		return false, err
	}
	return bytes.Equal(add.Bytes(), sigPublicKey), nil
}

// TODO
//func (a Address) RecoverSignedMsg(PrivateKey) {
//
//}

//sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
//if err != nil {
//	log.Fatal(err)
//}
//
//matches := bytes.Equal(sigPublicKey, publicKeyBytes)
//fmt.Println(matches) // true
//
//sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
//if err != nil {
//	log.Fatal(err)
//}
//
//sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
//matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
//fmt.Println(matches) // true
//
//signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
//verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
//fmt.Println(verified) // true
