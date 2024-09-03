package blockchain

import (
	"fmt"
	"time"

	evm "github.com/snzpool/balance_monitor/pkg/blockchain/evm"
	startknet "github.com/snzpool/balance_monitor/pkg/blockchain/starknet"
	common "github.com/snzpool/balance_monitor/pkg/common"
)

var gEVMList = []string{
	"eth", "ethereum", "bsc", "matic", "polygon", "heco", "ftm", "fatom", "arb", "arbitrum", "xdai", "avax", "avalanche", "harmony", "one", "metis", "evm",
}

var gStarknetList = []string{"starknet", "starknet_eth", "starknet_strk"}

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
	} else {
		fmt.Printf("%s is not supported now. Please contact administrator to add it\n", network)
		return -1
	}
}

func GetBalance(urlStr string, network string, address string) float64 {

	if common.InStringList(network, gEVMList) {
		result := evm.GetBalanceGo(urlStr, address)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = evm.GetBalanceGo(urlStr, address)
		}
		return result
	} else if common.InStringList(network, gStarknetList) {
		if network == "starknet_strk" {
			result := startknet.GetBalanceSTRK(urlStr, address)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = startknet.GetBalanceSTRK(urlStr, address)
			}
			return result
		} else {
			result := startknet.GetBalanceETH(urlStr, address)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = startknet.GetBalanceETH(urlStr, address)
			}
			return result
		}
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
