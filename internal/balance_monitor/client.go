package balance_monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	blockchain "github.com/snzpool/balance_monitor/pkg/blockchain"
	common "github.com/snzpool/balance_monitor/pkg/common"
)

func RunBalanceCheck() {
	InitPrometheus()
	interval := time.Duration(gConf.Frequency) * time.Second
	limiters := NewEndpointLimiters(EffectiveRpcRPS())
	ctx := context.Background()

	for {
		common.PrintCurTime()

		var wg sync.WaitGroup
		for i := 0; i < len(gConf.NetworkList); i++ {
			wg.Add(1)
			go func(oneNetwork Network) {
				defer wg.Done()
				checkNetwork(ctx, limiters, oneNetwork)
			}(gConf.NetworkList[i])
		}
		wg.Wait()

		time.Sleep(interval)
	}
}

func checkNetwork(ctx context.Context, limiters *EndpointLimiters, oneNetwork Network) {
	typ := oneNetwork.Type
	networkName := oneNetwork.Network
	fmt.Printf("type: %s, network: %s\n", typ, networkName)

	endpoints := oneNetwork.Endpoints
	if len(endpoints) == 0 {
		fmt.Printf("network %s has no endpoints\n", networkName)
		return
	}

	blockNumbers := make([]int, len(endpoints))
	var heightWg sync.WaitGroup
	for j := 0; j < len(endpoints); j++ {
		heightWg.Add(1)
		go func(idx int) {
			defer heightWg.Done()
			url := endpoints[idx]
			if err := limiters.Wait(ctx, url); err != nil {
				fmt.Printf("rate limit wait failed for %s: %v\n", url, err)
				blockNumbers[idx] = -1
				return
			}
			blockNumbers[idx] = int(blockchain.GetBlockHeight(url, typ))
			fmt.Printf("- %s: %d\n", url, blockNumbers[idx])
		}(j)
	}
	heightWg.Wait()

	selectIndex := 0
	for j := 1; j < len(endpoints); j++ {
		if blockNumbers[j] > blockNumbers[selectIndex] {
			selectIndex = j
		}
	}
	selectedUrl := endpoints[selectIndex]
	fmt.Printf("selectedUrl: %s\n", selectedUrl)

	var addrWg sync.WaitGroup
	for j := 0; j < len(oneNetwork.AddressList); j++ {
		addrWg.Add(1)
		go func(oneAddressInfo AddressInfo) {
			defer addrWg.Done()
			if err := limiters.Wait(ctx, selectedUrl); err != nil {
				fmt.Printf("rate limit wait failed for %s: %v\n", selectedUrl, err)
				updateMetrics(oneAddressInfo, networkName, oneNetwork.TokenAddress, oneNetwork.Symbol, -1)
				return
			}
			balance := blockchain.GetBalance(
				selectedUrl,
				typ,
				oneAddressInfo.Address,
				oneNetwork.TokenAddress,
				oneNetwork.TokenDecimals,
			)
			updateMetrics(oneAddressInfo, networkName, oneNetwork.TokenAddress, oneNetwork.Symbol, balance)
		}(oneNetwork.AddressList[j])
	}
	addrWg.Wait()
}

func updateMetrics(oneAddressInfo AddressInfo, networkName, tokenAddress, symbol string, balance float64) {
	address := oneAddressInfo.Address
	label := oneAddressInfo.Label
	infoThreshold := oneAddressInfo.InfoThreshold
	warnThreshold := oneAddressInfo.WarnThreshold

	balance_monitor_address_balance.WithLabelValues(label, networkName, address).Set(balance)
	recordSnapshot(networkName, label, address, tokenAddress, symbol, balance)

	var lowFlag float64
	var warnFlag float64
	var rpcBadFlag float64
	if balance < 0 {
		rpcBadFlag = 1
	} else if balance < infoThreshold && balance > warnThreshold {
		lowFlag = 1
	} else if balance < warnThreshold {
		warnFlag = 1
	}

	balance_monitor_rpc_bad.WithLabelValues(label, networkName).Set(rpcBadFlag)
	balance_monitor_balance_low.WithLabelValues(label, networkName, address).Set(lowFlag)
	balance_monitor_balance_empty.WithLabelValues(label, networkName, address).Set(warnFlag)

	fmt.Printf("%s, %s, %s, %f, %f, %f, %f\n", networkName, address, label, balance, rpcBadFlag, lowFlag, warnFlag)
}
