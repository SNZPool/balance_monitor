package balance_monitor

import (
	"encoding/json"
	"errors"
	"fmt"

	common "github.com/snzpool/balance_monitor/pkg/common"
)

//
type AddressInfo struct {
	Address       string  `json:"address"`
	Label         string  `json:"label"`
	InfoThreshold float64 `json:"infoThreshold"`
	WarnThreshold float64 `json:"warnThreshold"`
}

type Network struct {
	Network       string        `json:"network"`
	Endpoints     []string      `json:"endpoints"`
	TokenAddress  string        `json:"tokenAddress,omitempty"`  // empty = default asset for this network
	TokenDecimals *int          `json:"tokenDecimals,omitempty"` // optional override; nil = read decimals on-chain
	AddressList   []AddressInfo `json:"addressList"`
}

type Conf struct {
	Frequency   int       `json:"frequency"`
	MetricPort  int       `json:"metricPort"`
	RpcRPS      float64   `json:"rpcRPS"` // per-endpoint requests/sec; default 5 when <= 0
	NetworkList []Network `json:"info"`
}

const defaultRpcRPS = 5.0

//
var gConf Conf

// EffectiveRpcRPS returns the configured per-endpoint RPS, or the default when unset/invalid.
func EffectiveRpcRPS() float64 {
	if gConf.RpcRPS <= 0 {
		return defaultRpcRPS
	}
	return gConf.RpcRPS
}
//
func InitConfig(configPath string) error {

	err := ReadConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	PrintConfig()
	return nil
}

//
func ReadConfig(configPath string) error {
	// json string -> struct
	fmt.Printf("reading file %s\n", configPath)
	data, err := common.ReadFile(configPath)
	err = json.Unmarshal(data, &gConf)
	if err != nil {
		info := fmt.Sprintf("json.Unmarshal failed, err:%v\n", err)
		return errors.New(info)
	}

	return nil
}

//
func PrintConfig() {
	fmt.Println(gConf.Frequency)
	fmt.Println(gConf.MetricPort)
	fmt.Printf("rpcRPS: %g\n", EffectiveRpcRPS())
	for i := 0; i < len(gConf.NetworkList); i++ {
		fmt.Println(gConf.NetworkList[i].Network)
		if gConf.NetworkList[i].TokenAddress != "" {
			fmt.Printf("tokenAddress: %s\n", gConf.NetworkList[i].TokenAddress)
		}
		if gConf.NetworkList[i].TokenDecimals != nil {
			fmt.Printf("tokenDecimals: %d\n", *gConf.NetworkList[i].TokenDecimals)
		}
		for j := 0; j < len(gConf.NetworkList[i].Endpoints); j++ {
			fmt.Println(gConf.NetworkList[i].Endpoints[j])
		}
		for j := 0; j < len(gConf.NetworkList[i].AddressList); j++ {
			fmt.Println(gConf.NetworkList[i].AddressList[j].Address)
			fmt.Println(gConf.NetworkList[i].AddressList[j].Label)
			fmt.Println(gConf.NetworkList[i].AddressList[j].InfoThreshold)
			fmt.Println(gConf.NetworkList[i].AddressList[j].WarnThreshold)
		}
	}
}
