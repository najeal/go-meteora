package main

import (
	"context"
	"log"
	"os"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/joho/godotenv"
	"github.com/najeal/meteora-bot/positionmanager"
	"github.com/najeal/meteora-bot/pricechecker"
	"github.com/najeal/meteora-bot/store"
)

const (
	SolanaRPCEnvVar     = "SOLANA_RPC_ENDPOINT"
	MeteoraWalletEnvVar = "METEORA_WALLET"
	MeteoraProgramID    = "LBUZKhRxPF3XUpBCjp4YzTKgLccjZhTSDM9YuVaPwxo"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	solanaRpcEndpoint := os.Getenv(SolanaRPCEnvVar)
	walletAddress := os.Getenv(MeteoraWalletEnvVar)

	activeBinManager := pricechecker.NewActiveBinManager()
	st := store.NewPositionStore()
	positionManager := positionmanager.New(st)
	positionManager.Run(ctx, positionmanager.Config{
		SolanaRPCEndpoint: solanaRpcEndpoint,
		MeteoraProgramID:  MeteoraProgramID,
		WalletAddress:     common.PublicKeyFromString(walletAddress),
	}, st, activeBinManager)
}
