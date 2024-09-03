package starknet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "os"
	"strconv"
	"time"

	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

var (
	name                string = "mainnet"
	ethMainnetContract  string = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	strkMainnetContract string = "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"
	getBalanceMethod    string = "balanceOf"
	getDecimalsMethod   string = "decimals"
)

//
func GetBlockHeight(urlStr string) (resultNumber int64) {

	//
	jsonStr := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"starknet_blockNumber",
		"params":[],
		"id":0}`)

	//
	rest, err := sendJsonPOSTRequest(urlStr, jsonStr)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	//
	s := strconv.Itoa(rest.Result)
	num10, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return num10
}

//
func sendJsonPOSTRequest(urlStr string, jsonStr string) (rest Rest, state error) {

	//
	var reqMsg ReqMsg
	err := json.Unmarshal([]byte(jsonStr), &reqMsg)
	if err != nil {
		state = err
		return
	}
	byteData, _ := json.Marshal(&reqMsg)

	reader := bytes.NewReader([]byte(byteData))
	req, err := http.NewRequest("POST", urlStr, reader)
	if err != nil {
		state = err
		return
	}
	req.Header.Add("Content-Type", "application/json")

	//
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		state = err
		return
	}
	defer resp.Body.Close()

	//
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		state = err
		return
	}
	json.Unmarshal(respBytes, &rest)
	return
}

func GetBalanceETH(urlStr string, ethAddress string) float64 {
	return GetBalance(urlStr, ethAddress, ethMainnetContract)
}

func GetBalanceSTRK(urlStr string, ethAddress string) float64 {
	return GetBalance(urlStr, ethAddress, strkMainnetContract)
}

func GetBalance(urlStr string, ethAddress string, erc20ContractAddress string) float64 {
	clientv02, err := rpc.NewProvider(urlStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error dialing the RPC provider: %s", err))
	}

	// fmt.Println("Established connection with the client")

	tokenAddressInFelt, err := utils.HexToFelt(erc20ContractAddress)
	if err != nil {
		fmt.Println("Failed to transform the token contract address, did you give the hex address?")
		return -1.0
	}

	accountAddressInFelt, err := utils.HexToFelt(ethAddress)
	if err != nil {
		fmt.Println("Failed to transform the account address, did you give the hex address?")
		return -1.0
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getBalanceMethod),
		Calldata:           []*felt.Felt{accountAddressInFelt},
	}

	// fmt.Println("Making balanceOf() request")
	callResp, rpcErr := clientv02.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		return -1.0
	}

	// Get token's decimals
	getDecimalsTx := rpc.FunctionCall{
		ContractAddress:    tokenAddressInFelt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(getDecimalsMethod),
	}
	getDecimalsResp, rpcErr := clientv02.Call(context.Background(), getDecimalsTx, rpc.BlockID{Tag: "latest"})
	if rpcErr != nil {
		return -1.0
	}

	floatValue := new(big.Float).SetInt(utils.FeltToBigInt(callResp[0]))
	floatValue.Quo(floatValue, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), utils.FeltToBigInt(getDecimalsResp[0]), nil)))

	value, _ := floatValue.Float64()

	return value
}
