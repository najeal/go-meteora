package requester

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/najeal/meteora-bot/rpc"
	"go.uber.org/zap"
)

const (
	waitTime = 1
)

type Config struct {
	SolanaRPCEndpoint string
	MeteoraProgramID  string
	WalletAddress     common.PublicKey
}

func RecurrentFetch(logger *zap.Logger, cfg Config, stop <-chan struct{}) <-chan []rpc.PositionAccount {
	httpClient := &http.Client{}
	stream := make(chan []rpc.PositionAccount)
	go func() {
		defer close(stream)
		for {
			select {
			case <-stop:
				fmt.Println("stopping Recurrent Wallet Position fetching")
				return
			case <-time.After(time.Minute * waitTime):
				positionAccounts, err := rpc.GetWalletPositions(httpClient, cfg.SolanaRPCEndpoint, cfg.MeteoraProgramID, cfg.WalletAddress.String(), "")
				if err != nil {
					logger.Error("failed to fetch wallet positions",
						zap.String("wallet-address", cfg.WalletAddress.String()),
						zap.Error(err),
					)
					continue
				}
				stream <- positionAccounts

			}
		}
	}()
	return stream
}

func RecurrentFetchPairAccount(solanaRPCEndpoint string, meteoraPoolAddress common.PublicKey, stop <-chan struct{}) <-chan rpc.LbPairAccount {
	httpClient := &http.Client{}
	stream := make(chan rpc.LbPairAccount)
	go func() {
		defer close(stream)
		for {
			select {
			case <-stop:
				return
			case <-time.After(time.Minute * waitTime):
				pairAccount, err := rpc.GetLbPairAccount(httpClient, solanaRPCEndpoint, meteoraPoolAddress.String())
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				stream <- pairAccount
			}
		}
	}()
	return stream
}
