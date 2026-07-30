package blockchain

import (
	"fmt"
	"time"

	evm "github.com/snzpool/balance_monitor/pkg/blockchain/evm"
	startknet "github.com/snzpool/balance_monitor/pkg/blockchain/starknet"
	tron "github.com/snzpool/balance_monitor/pkg/blockchain/tron"
)

// Supported type values for protocol routing.
const (
	TypeEVM          = "evm"
	TypeStarknet     = "starknet"
	TypeStarknetSTRK = "starknet_strk"
	TypeTron         = "tron"
)

func GetBlockHeight(urlStr string, typ string) int64 {
	switch typ {
	case TypeEVM:
		result := evm.GetBlockHeight(urlStr)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = evm.GetBlockHeight(urlStr)
		}
		return result
	case TypeStarknet, TypeStarknetSTRK:
		result := startknet.GetBlockHeight(urlStr)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = startknet.GetBlockHeight(urlStr)
		}
		return result
	case TypeTron:
		result := tron.GetBlockHeight(urlStr)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = tron.GetBlockHeight(urlStr)
		}
		return result
	default:
		fmt.Printf("type %q is not supported (use evm, starknet, starknet_strk, or tron)\n", typ)
		return -1
	}
}

// GetBalance fetches the balance for address using the given protocol type.
// When tokenAddress is empty, uses the type default asset (native for EVM/Tron;
// ETH for starknet, STRK for starknet_strk). When set, queries that ERC20 / SNIP-20 contract.
func GetBalance(urlStr string, typ string, address string, tokenAddress string, tokenDecimals *int) float64 {
	switch typ {
	case TypeEVM:
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
	case TypeStarknet, TypeStarknetSTRK:
		if tokenAddress != "" {
			result := startknet.GetBalance(urlStr, address, tokenAddress, tokenDecimals)
			if result < 0 {
				time.Sleep(time.Duration(interval) * time.Second)
				result = startknet.GetBalance(urlStr, address, tokenAddress, tokenDecimals)
			}
			return result
		}
		if typ == TypeStarknetSTRK {
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
	case TypeTron:
		if tokenAddress != "" {
			fmt.Printf("tokenAddress is not supported for type %s\n", typ)
			return -1
		}
		result := tron.GetBalanceGo(urlStr, address)
		if result < 0 {
			time.Sleep(time.Duration(interval) * time.Second)
			result = tron.GetBalanceGo(urlStr, address)
		}
		return result
	default:
		fmt.Printf("type %q is not supported (use evm, starknet, starknet_strk, or tron)\n", typ)
		return -1
	}
}

func GetGasCost(urlStr string, typ string, address string, startTime string, endTime string) float64 {
	if typ == TypeEVM {
	} else {
		fmt.Println("waiting")
	}
	return 0
}
