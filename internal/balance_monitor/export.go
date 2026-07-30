package balance_monitor

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sync"
)

// BalanceSnapshot is one row for on-demand CSV export.
type BalanceSnapshot struct {
	Network      string
	Label        string
	Address      string
	Balance      float64
	TokenAddress string
	Symbol       string
}

var (
	snapshotMu    sync.RWMutex
	snapshotByKey = make(map[string]BalanceSnapshot)
)

func snapshotKey(network, label, address string) string {
	return network + "\x00" + label + "\x00" + address
}

func recordSnapshot(network, label, address, tokenAddress, symbol string, balance float64) {
	snapshotMu.Lock()
	defer snapshotMu.Unlock()
	snapshotByKey[snapshotKey(network, label, address)] = BalanceSnapshot{
		Network:      network,
		Label:        label,
		Address:      address,
		Balance:      balance,
		TokenAddress: tokenAddress,
		Symbol:       symbol,
	}
}

func snapshotRows() []BalanceSnapshot {
	snapshotMu.RLock()
	defer snapshotMu.RUnlock()
	rows := make([]BalanceSnapshot, 0, len(snapshotByKey))
	for _, row := range snapshotByKey {
		rows = append(rows, row)
	}
	return rows
}

func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows := snapshotRows()
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="balances.csv"`)

	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"network", "label", "address", "balance", "tokenAddress", "symbol"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, row := range rows {
		if err := cw.Write([]string{
			row.Network,
			row.Label,
			row.Address,
			fmt.Sprintf("%g", row.Balance),
			row.TokenAddress,
			row.Symbol,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
