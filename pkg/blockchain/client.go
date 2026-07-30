package blockchain

import (
	"fmt"
	"time"

	evm "github.com/snzpool/balance_monitor/pkg/blockchain/evm"
	startknet "github.com/snzpool/balance_monitor/pkg/blockchain/starknet"
	tron "github.com/snzpool/balance_monitor/pkg/blockchain/tron"
	common "github.com/snzpool/balance_monitor/pkg/common"
)

var gEVMList = []string{
	"eth", "ethereum", "bsc", "matic", "polygon", "heco", "ftm", "fatom", "arb", "arbitrum", "xdai", "avax", "avalanche", "harmony", "one", "metis", "evm", "tempo",
}

var gStarknetList = []string{"starknet", "starknet_eth", "starknet_strk"}

var gTronList = []string{"tron"}

func GetBlockHeight(urlStr string, network string) int64 {
	var result int64 = -1
	if common.InStringList(network, gEVMList) {
		result = evm.GetBlockHeight(urlStr)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = evm.GetBlockHeight(urlStr)
		}
		return result
	} else if common.InStringList(network, gStarknetList) {
		result := startknet.GetBlockHeight(urlStr)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = startknet.GetBlockHeight(urlStr)
		}
		return result
	} else if common.InStringList(network, gTronList) {
		result := tron.GetBlockHeight(urlStr)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = tron.GetBlockHeight(urlStr)
		}
		return result
	} else {
		fmt.Printf("%s is not supported now. Please contact administrator to add it\n", network)
		return -1
	}
}

// GetBalance fetches the balance for address on the given network.
// When tokenAddress is empty, uses the network default asset (native for EVM/Tron;
// ETH or STRK for Starknet aliases). When set, queries that ERC20 / SNIP-20 contract.
func GetBalance(urlStr string, network string, address string, tokenAddress string, tokenDecimals *int) float64 {

	if common.InStringList(network, gEVMList) {
		var result float64
		if tokenAddress != "" {
			result = evm.GetERC20Balance(urlStr, tokenAddress, address, tokenDecimals)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = evm.GetERC20Balance(urlStr, tokenAddress, address, tokenDecimals)
			}
		} else {
			result = evm.GetBalanceGo(urlStr, address)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = evm.GetBalanceGo(urlStr, address)
			}
		}
		return result
	} else if common.InStringList(network, gStarknetList) {
		if tokenAddress != "" {
			result := startknet.GetBalance(urlStr, address, tokenAddress, tokenDecimals)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = startknet.GetBalance(urlStr, address, tokenAddress, tokenDecimals)
			}
			return result
		}
		if network == "starknet_strk" {
			result := startknet.GetBalanceSTRK(urlStr, address)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = startknet.GetBalanceSTRK(urlStr, address)
			}
			return result
		}
		result := startknet.GetBalanceETH(urlStr, address)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = startknet.GetBalanceETH(urlStr, address)
		}
		return result
	} else if common.InStringList(network, gTronList) {
		if tokenAddress != "" {
			fmt.Printf("tokenAddress is not supported for network %s\n", network)
			return -1
		}
		result := tron.GetBalanceGo(urlStr, address)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = tron.GetBalanceGo(urlStr, address)
		}
		return result
	} else {
		fmt.Printf("%s is not supported now. Please contact administrator to add it\n", network)
		return -1
	}
}

func GetGasCost(urlStr string, network string, address string, startTime string, endTime string) float64 {
	if common.InStringList(network, gEVMList) {

	} else {
		fmt.Println("waiting")
	}
	return 0
}
