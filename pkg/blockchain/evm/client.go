package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//
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

//
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

//
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

func GetBalanceGo(urlStr string, ethAddress string) float64 {

	client, err := ethclient.Dial(urlStr)
	if err != nil {
		log.Fatal(err)
	}

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
