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
	Network     string        `json:"network"`
	Endpoints   []string      `json:"endpoints"`
	AddressList []AddressInfo `json:"addressList"`
}

type Conf struct {
	Frequency   int       `json:"frequency"`
	MetricPort  int       `json:"metricPort"`
	NetworkList []Network `json:"info"`
}

//
var gConf Conf

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
	for i := 0; i < len(gConf.NetworkList); i++ {
		fmt.Println(gConf.NetworkList[i].Network)
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
