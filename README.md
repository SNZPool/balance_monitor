[![Go Build Check](https://github.com/SNZPool/balance_monitor/actions/workflows/build.yml/badge.svg)](https://github.com/SNZPool/balance_monitor/actions/workflows/build.yml)

# Balance Monitor

Poll wallet balances across multiple chains and expose Prometheus gauges for balance and low-balance thresholds.

Supported chains:

| Family | Default asset (no `tokenAddress`) | Optional `tokenAddress` |
|--------|-----------------------------------|-------------------------|
| EVM (aliases below) | Native gas token (`eth_getBalance` / `BalanceAt`) | Any ERC20 (`balanceOf` + `decimals`) |
| Starknet | ETH (`starknet` / `starknet_eth`) or STRK (`starknet_strk`) | Any SNIP-20-style contract |
| Tron | Native TRX | Not supported |

Old configs without `tokenAddress` / `tokenDecimals` keep working unchanged.

## Requirements

- Go `>= 1.23.0`

## Install and build

```bash
make install   # go mod tidy
make build     # ./bin/balance_monitor (host OS)
# make build_linux   # linux/amd64 binary
```

## Usage

```bash
./bin/balance_monitor -config /path/to/config.json
```

Flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `balanceCheckConfig.json` | Path to the config file |

The loader expects **JSON** content. The sample file is named `depolyments/config-sample.toml` for historical reasons; the content is JSON.

Quick local run against the sample:

```bash
make test
```

Prometheus metrics are served at `http://<host>:<metricPort>/metrics`.

On-demand CSV export of the latest polled balances (no continuous file write):

```bash
curl -o balances.csv "http://<host>:<metricPort>/export.csv"
```

Columns: `network`, `label`, `address`, `balance`, `tokenAddress` (empty when monitoring the network default / native asset). Data reflects the last completed poll for each address; before the first cycle finishes the file may be empty or partial.

## Configuration

Full example: [`depolyments/config-sample.toml`](./depolyments/config-sample.toml).

### Top-level fields

| Field | Type | Description |
|-------|------|-------------|
| `frequency` | int | Poll interval in seconds (starts after the previous cycle finishes) |
| `metricPort` | int | HTTP port for `/metrics` |
| `rpcRPS` | float | Optional. Max requests per second **per endpoint URL** (default `5` when omitted or `<= 0`). Protects public RPCs from bursts |
| `info` | array | Per-chain (or per-asset) monitor groups |

### Network group (`info[]`)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `network` | string | yes | Chain selector (see [Network values](#network-values)) |
| `endpoints` | string[] | yes | RPC URLs; the one with the **highest block height** is used each cycle |
| `tokenAddress` | string | no | ERC20 / SNIP-20 contract. Empty = network default asset |
| `tokenDecimals` | int | no | Fallback decimals if on-chain `decimals()` fails |
| `addressList` | array | yes | Addresses to monitor |

### Address entry (`addressList[]`)

| Field | Type | Description |
|-------|------|-------------|
| `address` | string | Wallet / account address |
| `label` | string | Prometheus `name` label (use a unique, readable value) |
| `infoThreshold` | float | Soft low threshold → `balance_low = 1` |
| `warnThreshold` | float | Critical low threshold → `balance_empty = 1` |

Prefer `infoThreshold > warnThreshold`. Threshold flags are mutually exclusive: critical sets only `balance_empty`, not both.

### Network values

**EVM aliases** (native when `tokenAddress` is empty; ERC20 when set):

`eth`, `ethereum`, `bsc`, `matic`, `polygon`, `heco`, `ftm`, `fatom`, `arb`, `arbitrum`, `xdai`, `avax`, `avalanche`, `harmony`, `one`, `metis`, `evm`, `tempo`

**Starknet:**

| `network` | Default contract when `tokenAddress` is empty |
|-----------|-----------------------------------------------|
| `starknet`, `starknet_eth` | Starknet ETH |
| `starknet_strk` | STRK |

**Tron:** `tron` (native TRX only; `tokenAddress` is rejected).

### `tokenAddress` behavior

| `network` | `tokenAddress` | Behavior |
|-----------|----------------|----------|
| EVM alias | empty | Native balance |
| EVM alias | set | ERC20 `balanceOf` / `decimals` |
| `starknet` / `starknet_eth` | empty | Default ETH contract |
| `starknet_strk` | empty | Default STRK contract |
| any `starknet*` | set | Configured contract |
| `tron` | set | Error (`rpc_bad`) |

To monitor **two assets on the same chain**, use two `info` entries (for example native + ERC20, or Starknet ETH + STRK). Distinguishing assets in metrics relies on `network` and/or `label` — there is no `token` metric label.

### Examples

**Tempo pathUSD (gas token as ERC20):**

```json
{
  "network": "tempo",
  "endpoints": ["https://tempo-mainnet.drpc.org"],
  "tokenAddress": "0x20C0000000000000000000000000000000000000",
  "addressList": [
    {
      "address": "0xYourAddress",
      "label": "tempo_ops_pathusd",
      "infoThreshold": 100,
      "warnThreshold": 20
    }
  ]
}
```

If `decimals()` fails on a non-standard contract, set `"tokenDecimals": 6` (pathUSD uses 6 decimals on Tempo mainnet).

**Starknet STRK via explicit contract** (equivalent default to `network: "starknet_strk"` with empty `tokenAddress`):

```json
{
  "network": "starknet",
  "endpoints": ["https://starknet-mainnet.public.blastapi.io"],
  "tokenAddress": "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d",
  "addressList": [
    {
      "address": "0xYourStarknetAccount",
      "label": "starknet_strk_user1_account",
      "infoThreshold": 0.002,
      "warnThreshold": 0.001
    }
  ]
}
```

**Backward compatibility:** keep existing `network: "starknet_strk"` entries without `tokenAddress` so Prometheus `network` labels and alerts do not change. Migrating to `network: "starknet"` + `tokenAddress` changes the `network` label and may require updating dashboards/alerts.

**EVM native** (unchanged from older configs):

```json
{
  "network": "evm",
  "endpoints": [
    "https://andromeda.metis.io/?owner=1088",
    "https://metis-mainnet.public.blastapi.io"
  ],
  "addressList": [
    {
      "address": "0xYourAddress",
      "label": "metis_user1_account",
      "infoThreshold": 0.05,
      "warnThreshold": 0.02
    }
  ]
}
```

## How it works

Each poll cycle:

1. Run all `info` groups **in parallel** (different chains do not wait on each other).
2. Within a group, probe all `endpoints` for block height **in parallel**, then pick the highest.
3. Fetch all addresses in that group **in parallel** against the selected RPC.
4. Compare against thresholds and update Prometheus gauges.
5. When the whole cycle finishes, sleep `frequency` seconds (cycles do not overlap).

**Rate limiting:** every height/balance request waits on a per-endpoint-URL limiter (`rpcRPS`, default 5, burst 1). Different URLs are independent; the same URL shared by multiple groups shares one limiter. Failed RPC calls are still retried once after a short delay inside the blockchain client.

On fetch failure the balance is treated as `-1` and `rpc_bad` is set to `1`.

## Monitor metrics

Exported names use subsystem prefix `sdk_` (for example `sdk_balance_monitor_address_balance`).

| Metric | Labels | Description |
|--------|--------|-------------|
| `sdk_balance_monitor_address_balance` | `name`, `network`, `address` | Normalized balance (float) |
| `sdk_balance_monitor_rpc_bad` | `name`, `network` | `1` if the last fetch failed |
| `sdk_balance_monitor_balance_low` | `name`, `network`, `address` | `1` if `warnThreshold < balance < infoThreshold` |
| `sdk_balance_monitor_balance_empty` | `name`, `network`, `address` | `1` if `balance < warnThreshold` |

Sample series:

```
sdk_balance_monitor_address_balance{address="0x...",name="metis_user1_account",network="evm"} 0.335
sdk_balance_monitor_balance_low{address="0x...",name="metis_user1_account",network="evm"} 0
sdk_balance_monitor_balance_empty{address="0x...",name="metis_user1_account",network="evm"} 0
sdk_balance_monitor_rpc_bad{name="metis_user1_account",network="evm"} 0
```
