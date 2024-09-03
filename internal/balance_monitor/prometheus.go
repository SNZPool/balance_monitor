package balance_monitor

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	balance_monitor_address_balance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "sdk",
			Name:      "balance_monitor_address_balance",
			Help:      "balance_monitor_address_balance",
		}, []string{"name", "network", "address"})
	balance_monitor_rpc_bad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "sdk",
			Name:      "balance_monitor_rpc_bad",
			Help:      "balance_monitor_rpc_bad",
		}, []string{"name", "network"})
	balance_monitor_balance_low = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "sdk",
			Name:      "balance_monitor_balance_low",
			Help:      "balance_monitor_balance_low",
		}, []string{"name", "network", "address"})
	balance_monitor_balance_empty = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "sdk",
			Name:      "balance_monitor_balance_empty",
			Help:      "balance_monitor_balance_empty",
		}, []string{"name", "network", "address"})
)

func InitPrometheus() {
	//
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(balance_monitor_address_balance)
	prometheus.MustRegister(balance_monitor_rpc_bad)
	prometheus.MustRegister(balance_monitor_balance_low)
	prometheus.MustRegister(balance_monitor_balance_empty)

	//
	StartServer(int64(gConf.MetricPort))
}

func StartServer(port int64) (*http.Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}
	go func() {
		fmt.Println("Start the server.")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Httpserver: ListenAndServe() error: %s\n", err)
		}
	}()
	return srv, nil
}
