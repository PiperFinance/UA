package schemas

import "github.com/ethereum/go-ethereum/common"

type UserBalance struct {
	TokenStr  string         `bson:"tokenStr" json:"tokenStr"`
	UserStr   string         `bson:"userStr" json:"userStr"`
	User      common.Address `bson:"user" json:"user"`
	Token     common.Address `bson:"token" json:"token"`
	TokenId   TokenId        `bson:"token_id" json:"token_id"`
	TrxCount  uint64         `bson:"count" json:"count"`
	ChangedAt uint64         `bson:"c_t" json:"c_t"`
	StartedAt uint64         `bson:"s_t" json:"s_t"`
	Balance   string         `bson:"bal" json:"bal"`
}
