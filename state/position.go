package state

import (
	"math/big"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/rpc"
)

const (
	MAX_BIN_PER_POSITION = 70
	NUM_REWARDS          = 2
)

// NewPositionsWithAccountPubkey returns the datas contained into the account + the account pubkey
func NewPositionsWithAccountPubkey(positionAccounts []rpc.PositionAccount) ([]PositionWithAccountPubkey, error) {
	res := make([]PositionWithAccountPubkey, 0, len(positionAccounts))
	for _, positionAccount := range positionAccounts {
		position, err := NewPositionWithAccountPubkey(positionAccount)
		if err != nil {
			return nil, err
		}
		res = append(res, position)
	}
	return res, nil
}

// NewPositionWithAccountPubkey returns the data contained into the account + the account pubkey
func NewPositionWithAccountPubkey(positionAccount rpc.PositionAccount) (PositionWithAccountPubkey, error) {
	position, err := FromPositionAccount(positionAccount)
	if err != nil {
		return PositionWithAccountPubkey{}, err
	}
	return PositionWithAccountPubkey{
		Position:      position,
		AccountPubkey: positionAccount.Pubkey,
	}, nil
}

// FromPositionAccount returns the data contained into the account
func FromPositionAccount(positionAccount rpc.PositionAccount) (Position, error) {
	position, err := DeserializeData[Position](positionAccount.Account.Data[0])
	if err != nil {
		return Position{}, err
	}
	return position, nil
}

// FromPositionAccounts returns the datas contained into the account
func FromPositionAccounts(positionAccounts []rpc.PositionAccount) ([]Position, error) {
	positions := make([]Position, 0, len(positionAccounts))
	for _, positionAccount := range positionAccounts {
		position, err := FromPositionAccount(positionAccount)
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	return positions, nil
}

type PositionWithAccountPubkey struct {
	Position      Position
	AccountPubkey common.PublicKey
}

type UserRewardInfo struct {
	Reward_per_token_completes [NUM_REWARDS]big.Int
	Reward_pendings            [NUM_REWARDS]uint64
}

type FeeInfo struct {
	Fee_x_per_token_complete big.Int
	Fee_y_per_token_complete big.Int
	Fee_x_pending            uint64
	Fee_y_pending            uint64
}

type Position struct {
	Prefix                     [8]byte
	Lb_pair                    common.PublicKey
	Owner                      common.PublicKey
	Liquidity_shares           [MAX_BIN_PER_POSITION]big.Int
	Reward_infos               [MAX_BIN_PER_POSITION]UserRewardInfo
	Fee_infos                  [MAX_BIN_PER_POSITION]FeeInfo
	Lower_bin_id               int32
	Upper_bin_id               int32
	Last_updated_at            int64
	Total_claimed_fee_x_amount uint64
	Total_claimed_fee_y_amount uint64
	Total_claimed_rewards      [2]uint64
	Reserved                   [160]uint8
}
