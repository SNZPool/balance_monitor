package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetBlockHeight(urlStr string) (resultNumber int64) {

	//
	jsonStr := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_blockNumber",
		"params":[],
		"id":1}`)

	//
	rest, err := sendJsonPOSTRequest(urlStr, jsonStr)
	if err != nil {
		return -1
	}

	//
	num10, err := strconv.ParseInt(rest.Result, 0, 64)
	if err != nil {
		// fmt.Println(err)
		return -1
	}

	return num10
}

func GetBalance(urlStr string, ethAddress string) float64 {

	//
	jsonStr := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getBalance",
		"params":["%s", "latest"],
		"id":1}`, ethAddress)

	//
	rest, err := sendJsonPOSTRequest(urlStr, jsonStr)
	if err != nil {
		return -1
	}

	//
	num10, err := strconv.ParseInt(rest.Result, 0, 64)
	if err != nil {
		// fmt.Println(err)
		return -1
	}
	balance := float64(num10) / 1000000000000000000

	return balance
}

func GetTransCount(urlStr string, ethAddress string) int64 {

	//
	jsonStr := fmt.Sprintf(`{
		"jsonrpc":"2.0",
		"method":"eth_getTransactionCount",
		"params":["%s", "latest"],
		"id":1}`, ethAddress)

	//
	rest, err := sendJsonPOSTRequest(urlStr, jsonStr)
	if err != nil {
		// fmt.Println("Error: sendJsonPOSTRequest can't get results.")
		// fmt.Println(err)
		return -1
	}

	//
	num10, err := strconv.ParseInt(rest.Result, 0, 64)
	if err != nil {
		// fmt.Println(err)
		return -1
	}

	return num10
}

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

func GetBalanceGo(urlStr string, ethAddress string) float64 {

	client, err := ethclient.Dial(urlStr)
	if err != nil {
		fmt.Printf("ethclient.Dial failed: %v\n", err)
		return -1
	}
	defer client.Close()

	account := common.HexToAddress(ethAddress)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return -1
	}
	//fmt.Println(balance)

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValueBig := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	ethValue, _ := ethValueBig.Float64()
	//fmt.Println(ethValue)

	return ethValue
}

// ERC20 method selectors (first 4 bytes of keccak256)
var (
	balanceOfSelector = []byte{0x70, 0xa0, 0x82, 0x31} // balanceOf(address)
	decimalsSelector  = []byte{0x31, 0x3c, 0xe5, 0x67} // decimals()
)

// GetERC20Balance returns the ERC20 token balance for owner, normalized by decimals.
// If decimalsOverride is non-nil, it is used instead of (or as fallback for) the on-chain decimals() call.
func GetERC20Balance(urlStr string, tokenAddress string, ownerAddress string, decimalsOverride *int) float64 {
	client, err := ethclient.Dial(urlStr)
	if err != nil {
		fmt.Printf("ethclient.Dial failed: %v\n", err)
		return -1
	}
	defer client.Close()

	token := common.HexToAddress(tokenAddress)
	owner := common.HexToAddress(ownerAddress)

	// balanceOf(address)
	balanceData := make([]byte, 4+32)
	copy(balanceData[0:4], balanceOfSelector)
	copy(balanceData[4+12:], owner.Bytes()) // address left-padded to 32 bytes

	balanceMsg := ethereum.CallMsg{To: &token, Data: balanceData}
	balanceRaw, err := client.CallContract(context.Background(), balanceMsg, nil)
	if err != nil || len(balanceRaw) < 32 {
		fmt.Printf("ERC20 balanceOf failed: %v\n", err)
		return -1
	}
	balance := new(big.Int).SetBytes(balanceRaw)

	decimals := -1
	decimalsData := make([]byte, 4)
	copy(decimalsData[0:4], decimalsSelector)
	decimalsMsg := ethereum.CallMsg{To: &token, Data: decimalsData}
	decimalsRaw, err := client.CallContract(context.Background(), decimalsMsg, nil)
	if err == nil && len(decimalsRaw) > 0 {
		decimals = int(new(big.Int).SetBytes(decimalsRaw).Int64())
	} else if decimalsOverride != nil {
		decimals = *decimalsOverride
	} else {
		fmt.Printf("ERC20 decimals() failed: %v\n", err)
		return -1
	}
	if decimals < 0 || decimals > 77 {
		fmt.Printf("invalid ERC20 decimals: %d\n", decimals)
		return -1
	}

	fbalance := new(big.Float).SetInt(balance)
	divisor := new(big.Float).SetFloat64(math.Pow10(decimals))
	valueBig := new(big.Float).Quo(fbalance, divisor)
	value, _ := valueBig.Float64()
	return value
}
