package pricechecker

import (
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/state"
)

type LbPairActiveBin struct {
	Address     common.PublicKey
	ActiveBinID int32
}

func Check(poolID common.PublicKey, positions []state.PositionWithAccountPubkey, activeBinID int32) {
	for _, position := range positions {
		switch {
		case activeBinID < position.Position.Lower_bin_id:
			fmt.Printf("poolID %s : Out from the Bottom\n", poolID.String())
		case activeBinID > position.Position.Upper_bin_id:
			fmt.Printf("poolID %s : Out from the Top\n", poolID.String())
		default:
			fmt.Printf("poolID %s : lowerBin %v - activeBin %v - upperBin %v\n", poolID.String(), position.Position.Lower_bin_id, activeBinID, position.Position.Upper_bin_id)
		}
	}
}
