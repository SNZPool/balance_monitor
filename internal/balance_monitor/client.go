package balance_monitor

import (
	"fmt"
	"time"

	blockchain "github.com/snzpool/balance_monitor/pkg/blockchain"
	common "github.com/snzpool/balance_monitor/pkg/common"
)

func RunBalanceCheck() {

	//
	InitPrometheus()
	var interval time.Duration = time.Duration(gConf.Frequency) * time.Second

	var lowFlag float64 = 0
	var warnFlag float64 = 0
	var rpcBadFlag float64 = 0

	bLoop := true
	for bLoop {
		common.PrintCurTime()

		var selectIndex int = 0
		var selectedUrl string = ""
		for i := 0; i < len(gConf.NetworkList); i++ {
			oneNetwork := gConf.NetworkList[i]
			networkName := oneNetwork.Network
			fmt.Printf("network: %s\n", networkName)

			// select the rpc with highest block
			var blockNumberList []int
			for j := 0; j < len(oneNetwork.Endpoints); j++ {
				blockNumber := int(blockchain.GetBlockHeight(oneNetwork.Endpoints[j], networkName))
				fmt.Printf("- %s: %d\n", oneNetwork.Endpoints[j], blockNumber)
				blockNumberList = append(blockNumberList, blockNumber)
			}
			selectIndex = 0
			selectedUrl = ""
			for j := 0; j < len(oneNetwork.Endpoints); j++ {
				if j == 0 {
					selectIndex = 0
				} else {
					if blockNumberList[j] > blockNumberList[selectIndex] {
						selectIndex = j
					}
				}
			}
			selectedUrl = oneNetwork.Endpoints[selectIndex]
			fmt.Printf("selectedUrl: %s\n", selectedUrl)

			// loop
			for j := 0; j < len(oneNetwork.AddressList); j++ {
				oneAddressInfo := oneNetwork.AddressList[j]
				address := oneAddressInfo.Address
				label := oneAddressInfo.Label
				infoThreshold := oneAddressInfo.InfoThreshold
				warnThreshold := oneAddressInfo.WarnThreshold

				balance := blockchain.GetBalance(selectedUrl, networkName, address)
				balance_monitor_address_balance.WithLabelValues(label, networkName, address).Set(balance)

				lowFlag = 0
				warnFlag = 0
				rpcBadFlag = 0
				if balance < 0 {
					rpcBadFlag = 1
				} else if balance < infoThreshold && balance > warnThreshold {
					lowFlag = 1
				} else if balance < warnThreshold {
					warnFlag = 1
				}

				//
				balance_monitor_rpc_bad.WithLabelValues(label, networkName).Set(rpcBadFlag)
				balance_monitor_balance_low.WithLabelValues(label, networkName, address).Set(lowFlag)
				balance_monitor_balance_empty.WithLabelValues(label, networkName, address).Set(warnFlag)

				fmt.Printf("%s, %s, %s, %f, %f, %f, %f\n", networkName, address, label, balance, rpcBadFlag, lowFlag, warnFlag)
			}
		}

		time.Sleep(interval)
	}

}
