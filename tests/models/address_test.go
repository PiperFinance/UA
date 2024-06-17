package test

import (
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func FromPrivKey(privKey string) {

}

var (
	Addresses = []string{
		"fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19", //0x96216849c49358B10257cb55b28eA603c874b05E
		"f1cdbf9fc5dfc2dee0dd9b264520010062b0a86b60dc3ae672a4a8b245f2951d"} //0x0560aDB38A0C4828be88aeAE575F9ea5Acb549e8
)

func SignMsg(PrvKey string, Msg string) []byte {
	privateKey, _ := crypto.HexToECDSA(PrvKey)
	data := []byte(Msg)
	hash := crypto.Keccak256Hash(data)
	signature, _ := crypto.Sign(hash.Bytes(), privateKey)
	return signature
}

func TestVerifyMsg(t *testing.T) {

}
