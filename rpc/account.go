package rpc

import (
	"github.com/blocto/solana-go-sdk/common"
)

type Account struct {
	Data       [2]string `json:"data"`
	Executable bool      `json:"executable"`
	Lamports   uint64    `json:"lamports"`
	Owner      string    `json:"owner"`
	RentEpoch  uint64    `json:"rentEpoch"`
	Space      uint64    `json:"space"`
}

type PositionAccount struct {
	Pubkey  common.PublicKey `json:"pubkey"`
	Account Account          `json:"account"`
}

type LbPairAccount struct {
	Value Account
}
