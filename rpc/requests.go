package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	discriminatorPositionV2 = "LgkNAEYaVX3"
)

// GetLbPairAccount returns a LbPairAccount
func GetLbPairAccount(client HttpClient, solanaRPCEndpoint, meteoraPoolAddress string) (LbPairAccount, error) {
	// TODO replace by logger
	fmt.Println("GET LB PAIR", meteoraPoolAddress)
	reqBody := RPCRequest{
		JsonRPC: "2.0",
		Method:  "getAccountInfo",
		Params: []interface{}{
			meteoraPoolAddress,
			map[string]any{
				"commitment": "confirmed",
				"encoding":   "base64",
			},
		},
		ID: 1,
	}
	return sendRequest[LbPairAccount](client, reqBody, solanaRPCEndpoint)
}

// GetWalletPositions returns an array of PositionAccount
// It expects a meteoraPoolAddress to return one specific PositionAccount
// Without meteoraPoolAddress, the function returns all the PositionAccount of the wallet
func GetWalletPositions(client HttpClient, solanaRpcEndpoint, meteoraProgramID, walletAddress, meteoraPoolAddress string) ([]PositionAccount, error) {
	// TODO replace by logger
	fmt.Println("GET Wallet Positions", walletAddress)
	filters := []map[string]any{
		{
			"memcmp": map[string]any{
				"offset": 0,
				"bytes":  discriminatorPositionV2,
			},
		},
		{
			"memcmp": map[string]any{
				"offset": 40,
				"bytes":  walletAddress,
			},
		},
	}
	if meteoraPoolAddress != "" {
		filters = append(filters, map[string]any{
			"memcmp": map[string]any{
				"offset": 40,
				"bytes":  walletAddress,
			},
		})
	}
	reqBody := RPCRequest{
		JsonRPC: "2.0",
		Method:  "getProgramAccounts",
		Params: []interface{}{meteoraProgramID, map[string]any{
			"commitment": "confirmed",
			"encoding":   "base64",
			"filters":    filters,
		}},
		ID: 2,
	}

	return sendRequest[[]PositionAccount](client, reqBody, solanaRpcEndpoint)
}

func sendRequest[Result any](client HttpClient, rpcReq RPCRequest, solanaRpcEndpoint string) (out Result, err error) {
	// Encode the request into JSON
	jsonReq, err := json.Marshal(rpcReq)
	if err != nil {
		return out, err
	}
	req, err := http.NewRequest(http.MethodPost, solanaRpcEndpoint, bytes.NewBuffer(jsonReq))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	var rpcRes RPCResponse
	err = json.Unmarshal(body, &rpcRes)
	if err != nil {
		return out, err
	}

	err = json.Unmarshal(rpcRes.Result, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}
