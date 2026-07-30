[![Go Build Check](https://github.com/SNZPool/balance_monitor/actions/workflows/build.yml/badge.svg)](https://github.com/SNZPool/balance_monitor/actions/workflows/build.yml)

# Balance Monitor

Poll wallet balances across multiple chains and expose Prometheus gauges for balance and low-balance thresholds.

Supported protocols (`type`):

| `type` | Default asset (no `tokenAddress`) | Optional `tokenAddress` |
|--------|-----------------------------------|-------------------------|
| `evm` | Native gas token (`BalanceAt`) | Any ERC20 (`balanceOf` + `decimals`) |
| `starknet` | Starknet ETH | Any SNIP-20-style contract |
| `starknet_strk` | STRK | Any SNIP-20-style contract |
| `tron` | Native TRX | Not supported |

`network` is the **chain identity** used in Prometheus labels and CSV export (e.g. `eth`, `base`, `tempo`). It is independent of `type`, so you do not need to register new chain aliases in code.

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

The loader expects **JSON** content. See the sample at [`deployments/config-sample.json`](./deployments/config-sample.json).

Quick local run against the sample:

```bash
make test
```

Prometheus metrics are served at `http://<host>:<metricPort>/metrics`.

On-demand CSV export of the latest polled balances (no continuous file write):

```bash
curl -o balances.csv "http://<host>:<metricPort>/export.csv"
```

Columns: `network`, `label`, `address`, `balance`, `tokenAddress`, `symbol`. `tokenAddress` / `symbol` are empty when unset. Data reflects the last completed poll for each address; before the first cycle finishes the file may be empty or partial.

## Configuration

Full example: [`deployments/config-sample.json`](./deployments/config-sample.json).

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
| `type` | string | yes | Protocol: `evm`, `starknet`, `starknet_strk`, or `tron` |
| `network` | string | yes | Chain identity for metrics/export (e.g. `eth`, `base`, `tempo`). Free-form string |
| `endpoints` | string[] | yes | RPC URLs; the one with the **highest block height** is used each cycle |
| `symbol` | string | no | Gas / token symbol for CSV export and multi-chain stats (e.g. `ETH`, `pathUSD`). Display-only |
| `tokenAddress` | string | no | ERC20 / SNIP-20 contract. Empty = type default asset |
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

### `tokenAddress` behavior

| `type` | `tokenAddress` | Behavior |
|--------|----------------|----------|
| `evm` | empty | Native balance |
| `evm` | set | ERC20 `balanceOf` / `decimals` |
| `starknet` | empty | Default ETH contract |
| `starknet_strk` | empty | Default STRK contract |
| `starknet` / `starknet_strk` | set | Configured contract |
| `tron` | set | Error (`rpc_bad`) |

To monitor **two assets on the same chain**, use two `info` entries with the same `network` but different `type` / `tokenAddress` / `symbol` (and distinct `label`s).

### Examples

**Tempo pathUSD (gas token as ERC20):**

```json
{
  "type": "evm",
  "network": "tempo",
  "symbol": "pathUSD",
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

**Base (L2 using ETH as gas):**

```json
{
  "type": "evm",
  "network": "base",
  "symbol": "ETH",
  "endpoints": ["https://mainnet.base.org"],
  "addressList": [
    {
      "address": "0xYourAddress",
      "label": "chainlink_base_ocr",
      "infoThreshold": 0.05,
      "warnThreshold": 0.02
    }
  ]
}
```

**Starknet STRK:**

```json
{
  "type": "starknet_strk",
  "network": "starknet",
  "symbol": "STRK",
  "endpoints": ["https://starknet-mainnet.public.blastapi.io"],
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

**EVM native (Metis):**

```json
{
  "type": "evm",
  "network": "metis",
  "symbol": "METIS",
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

Exported names use subsystem prefix `sdk_` (for example `sdk_balance_monitor_address_balance`). The Prometheus `network` label is the config `network` field (chain identity), not `type`.

| Metric | Labels | Description |
|--------|--------|-------------|
| `sdk_balance_monitor_address_balance` | `name`, `network`, `address` | Normalized balance (float) |
| `sdk_balance_monitor_rpc_bad` | `name`, `network` | `1` if the last fetch failed |
| `sdk_balance_monitor_balance_low` | `name`, `network`, `address` | `1` if `warnThreshold < balance < infoThreshold` |
| `sdk_balance_monitor_balance_empty` | `name`, `network`, `address` | `1` if `balance < warnThreshold` |

Sample series:

```
sdk_balance_monitor_address_balance{address="0x...",name="metis_user1_account",network="metis"} 0.335
sdk_balance_monitor_balance_low{address="0x...",name="metis_user1_account",network="metis"} 0
sdk_balance_monitor_balance_empty{address="0x...",name="metis_user1_account",network="metis"} 0
sdk_balance_monitor_rpc_bad{name="metis_user1_account",network="metis"} 0
```
