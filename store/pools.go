package store

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/state"
)

// PositionStore stores a wallet positions
type PositionStore struct {
	// mapping pool address with the associated positions
	poolsPositions map[common.PublicKey][]state.PositionWithAccountPubkey
}

func NewPositionStore() *PositionStore {
	return &PositionStore{
		poolsPositions: map[common.PublicKey][]state.PositionWithAccountPubkey{},
	}
}

// GetPoolPositions returns the Positions associated to the poolAddress
func (x *PositionStore) GetPoolPositions(poolAddress common.PublicKey) ([]state.PositionWithAccountPubkey, bool) {
	positions, ok := x.poolsPositions[poolAddress]
	return positions, ok
}

// SetPositions store the Positions and removed the old ones
// It returns the list of removed and added Pools
func (x *PositionStore) SetPositions(positions []state.PositionWithAccountPubkey) (removed []common.PublicKey, added []common.PublicKey) {
	newstore := make(map[common.PublicKey][]state.PositionWithAccountPubkey)
	for _, position := range positions {
		var positionsToStore []state.PositionWithAccountPubkey
		if foundPositions, ok := newstore[position.Position.Lb_pair]; ok {
			positionsToStore = foundPositions
		}
		positionsToStore = append(positionsToStore, position)
		newstore[position.Position.Lb_pair] = positionsToStore
	}

	for lbpair := range newstore {
		if _, ok := x.poolsPositions[lbpair]; !ok {
			added = append(added, lbpair)
		}
	}

	for lbpair := range x.poolsPositions {
		if _, ok := newstore[lbpair]; !ok {
			removed = append(removed, lbpair)
		}
	}

	x.poolsPositions = newstore
	return removed, added
}
