package positionmanager

import (
	"context"
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/pricechecker"
	"github.com/najeal/meteora-bot/rpc"
	"github.com/najeal/meteora-bot/rpc/requester"
	"github.com/najeal/meteora-bot/state"
	"github.com/najeal/meteora-bot/store"
	"go.uber.org/zap"
)

type Storer interface {
	SetPositions([]state.PositionWithAccountPubkey) ([]common.PublicKey, []common.PublicKey)
	GetPoolPositions(poolAddress common.PublicKey) ([]state.PositionWithAccountPubkey, bool)
}

type PositionManager struct {
	PositionStore *store.PositionStore
}

func New(st *store.PositionStore) *PositionManager {
	return &PositionManager{
		PositionStore: st,
	}
}

type Config struct {
	SolanaRPCEndpoint string
	MeteoraProgramID  string
	WalletAddress     common.PublicKey
}

func (x *PositionManager) Run(ctx context.Context, cfg Config, httpClient *rpc.ClientLimiter, st Storer, binManager *pricechecker.ActiveBinManager) {
	recurrentFetchStop := make(chan struct{})
	defer close(recurrentFetchStop)
	paStream := requester.RecurrentFetch(zap.NewNop(), httpClient, requester.Config{
		SolanaRPCEndpoint: cfg.SolanaRPCEndpoint,
		MeteoraProgramID:  cfg.MeteoraProgramID,
		WalletAddress:     cfg.WalletAddress,
	}, recurrentFetchStop)
	binManagerStop := make(chan struct{})
	defer close(binManagerStop)
	updatesChan := make(chan pricechecker.TrackerUpdates)
	defer close(updatesChan)
	activeBinTracker := binManager.Run(cfg.SolanaRPCEndpoint, updatesChan, binManagerStop)
	for {
		select {
		case <-ctx.Done():
			// TODO replace by logger
			fmt.Println("stop position manager")
			return
		case positionAccounts := <-paStream:
			positions, err := state.NewPositionsWithAccountPubkey(positionAccounts)
			if err != nil {
				// TODO replace by logger
				fmt.Println("failed to transform position account into positions")
				continue
			}
			removedPools, addedPools := st.SetPositions(positions)
			updatesChan <- pricechecker.TrackerUpdates{
				Removed: removedPools,
				Added:   addedPools,
			}
		case activeBin := <-activeBinTracker:
			positions, ok := st.GetPoolPositions(activeBin.PoolAddress)
			if !ok {
				// TODO use logger
				continue
			}
			go pricechecker.Check(activeBin.PoolAddress, positions, activeBin.BinID)

		}
	}
}
