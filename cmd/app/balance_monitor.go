package main

import (
	"flag"

	balance_monitor "github.com/snzpool/balance_monitor/internal/balance_monitor"
)

var configPath = flag.String("config", "balanceCheckConfig.json", "the path of config")

func main() {
	flag.Parse()
	err := balance_monitor.InitConfig(*configPath)
	if err != nil {
		return
	}
	balance_monitor.RunBalanceCheck()
}
