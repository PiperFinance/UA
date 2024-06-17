package schemas

import "github.com/ethereum/go-ethereum/common"

type SyncUser struct {
	Hash common.Address `json:"hash"`
}
