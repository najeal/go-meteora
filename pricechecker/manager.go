package pricechecker

import (
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/rpc"
	"github.com/najeal/meteora-bot/rpc/requester"
	"github.com/najeal/meteora-bot/state"
)

type ActiveBinManager struct {
	binFetchers map[common.PublicKey]chan struct{}
	HttpClient  *rpc.ClientLimiter
}

func NewActiveBinManager(httpClient *rpc.ClientLimiter) *ActiveBinManager {
	return &ActiveBinManager{
		binFetchers: map[common.PublicKey]chan struct{}{},
		HttpClient:  httpClient,
	}
}

type TrackerUpdates struct {
	Removed []common.PublicKey
	Added   []common.PublicKey
}

type BinTacker struct {
	PoolAddress common.PublicKey
	BinID       int32
}

func (x *ActiveBinManager) Run(solanaRPCEndpoint string, updates <-chan TrackerUpdates, stop <-chan struct{}) chan BinTacker {
	binPriceStream := make(chan BinTacker)
	go func() {
		defer close(binPriceStream)
		for {
			select {
			case <-stop:
				// TODO replace by logger
				fmt.Println("stopping tracked bin recurrent fetchers")
				for _, stopChannel := range x.binFetchers {
					close(stopChannel)
				}
				// TODO replace by logger
				fmt.Println("stopping active bin manager")
				return
			case tracker := <-updates:
				for _, removed := range tracker.Removed {
					binChan, ok := x.binFetchers[removed]
					if !ok {
						continue
					}
					binChan <- struct{}{}
					delete(x.binFetchers, removed)
				}
				for _, added := range tracker.Added {
					binChan := make(chan struct{})
					binChanReader := requester.RecurrentFetchPairAccount(x.HttpClient, solanaRPCEndpoint, added, binChan)
					go func() {
						for pairAccount := range binChanReader {
							lbPair, err := state.FromLbPairAccount(pairAccount)
							if err != nil {
								fmt.Println("failed to get data from lb pair account, data:")
								fmt.Println(pairAccount.Value.Data[0])

								panic(err.Error())
							}
							binTracker := BinTacker{
								PoolAddress: added,
								BinID:       lbPair.ActiveID,
							}
							binPriceStream <- binTracker
						}
					}()
					x.binFetchers[added] = binChan

				}
			}
		}
	}()
	return binPriceStream
}
