package state

import (
	"math/big"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/rpc"
)

// FromLbPairAccount returns the data contained into the account
func FromLbPairAccount(pairAccount rpc.LbPairAccount) (LbPair, error) {
	lbPair, err := DeserializeData[LbPair](pairAccount.Value.Data[0])
	if err != nil {
		return LbPair{}, err
	}
	return lbPair, nil
}

type ProtocolFee struct {
	AmountX uint64
	AmountY uint64
}

type RewardInfo struct {
	Mint              common.PublicKey
	Vault             common.PublicKey
	Funder            common.PublicKey
	RewardDuration    uint64
	RewardDurationEnd uint64
	RewardRate        big.Int
	LastUpdateTime    uint64

	CumulativeSecondsWithEmptyLiquidityReward uint64
}

type LbPair struct {
	Prefix                  [8]byte
	Parameters              StaticParameters
	VParameters             VariableParameters
	BumpSeed                [1]uint8
	BinStepSeed             [2]uint8
	PairType                uint8
	ActiveID                int32
	BinStep                 uint16
	Status                  uint8
	RequireBaseFactorSeed   uint8
	BaseFactorSeed          [2]uint8
	ActivationType          uint8
	Padding0                uint8
	TokenXMint              common.PublicKey
	TokenYMint              common.PublicKey
	ReserveX                common.PublicKey
	ReserveY                common.PublicKey
	ProtocolFee             ProtocolFee
	Padding1                [32]uint8
	RewardInfos             [2]RewardInfo
	Orable                  common.PublicKey
	BinArrayBitmap          [16]uint64
	LastUpdatedAt           int64
	WhitelistedWallet       common.PublicKey
	PreActivationSwapWallet common.PublicKey
	BaseKey                 common.PublicKey
	ActivationPoint         uint64
	PreActivationDuration   uint64
	Padding2                [8]uint8
	LockDuration            uint64
	Creator                 common.PublicKey
	Reserved                [24]uint8
}
