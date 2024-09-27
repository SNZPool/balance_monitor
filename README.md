[![Go Build Check](https://github.com/SNZPool/balance_monitor/actions/workflows/build.yml/badge.svg)](https://github.com/SNZPool/balance_monitor/actions/workflows/build.yml)

# Introduction

This project is designed to retrieve the gas balance from multiple addresses across various blockchains. Currently, it supports EVM chains and the Starknet chain.

Additionally, it allows you to monitor whether an address's balance falls below specified thresholds. You can set up to two levels of thresholds for monitoring purposes.


# Install

At first, install `golang(>=1.23.0)`.

Then, install and build
```
make install
make build
```

The binary `balance_monitor` can be found in path `./bin`.

# Usage

You must create a config file and run it as below command.

```
./balance_monitor -config config.toml
```

# How to edit the config file

You can find a example config file at path `./depolyments/config-sample.toml`
```
{
  "frequency": 300,
  "metricPort": 8888,
  "info": [
    {
      "network": "evm",
      "endpoints": [
          "https://andromeda.metis.io/?owner=1088",
          "https://metis-mainnet.public.blastapi.io"
      ],
      "addressList": [
        {
          "address":"0x282976F47d7C1cE8B9f956241105d5741f203fA6",
          "label":"metis_user1_account",
          "infoThreshold": 0.05,
          "warnThreshold": 0.02
        },
        {
          "address":"0x81251B395506b4BfF35DB657751285c138BFB3F5",
          "label":"metis_user2_account",
          "infoThreshold": 0.1,
          "warnThreshold": 0.02
        }
      ]
    },
    {
      "network": "startknet",
      "endpoints": [
          "https://starknet-mainnet.public.blastapi.io",
          "https://free-rpc.nethermind.io/mainnet-juno"
      ],
      "addressList": [
        {
          "address":"0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",
          "label":"starknet_user1_account",
          "infoThreshold": 0.002,
          "warnThreshold": 0.001
        }
      ]
    },
    {
      "network": "startknet_strk",
      "endpoints": [
          "https://starknet-mainnet.public.blastapi.io",
          "https://free-rpc.nethermind.io/mainnet-juno"
      ],
      "addressList": [
        {
          "address":"0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",
          "label":"starknet_strk_user1_account",
          "infoThreshold": 0.002,
          "warnThreshold": 0.001
        }
      ]
    }
  ]
}
```

## Structure of toml config file

**Common Parameter**
- frequency: how long (seconds) the program to update the balance
- metricPort: the metric port for prometheus
- info: monitor info arrays, divided by blockchains

**Blockchain Parameter**
- network: support `evm`(Native Token), `starknet`(ETH) and `starknet_strk`(STRK)
- endpoins: you can enter multiple rpc's at the same time. Only the rpc with the highest height will be used
- addressList: you can enter multiple address for monitoring

**AddressList Parameter**
- address: the address
- label: used to distinguish address names by metrics label
- infoThreshold: when the balance of the address is below than this value, the monitor metric `balance_monitor_balance_low` is 1
- warnThreshold:  when the balance of the address is below than this value, the monitor metric `balance_monitor_balance_empty` is 1


## Monitor Metrics

| Metric Name | Description | Value |
| ---- | ---- | ---- |
| balance_monitor_address_balance | the gas balance in your address | float |
| balance_monitor_rpc_bad | is there an inaccessible rpc in the rpc you are using? | 0 or 1(yes) |
| balance_monitor_balance_low | whether the balance is below the infoThreshold | 0 or 1(yes) |
| balance_monitor_balance_empty | whether the balance is below the warnThreshold | 0 or 1(yes) |

**Sample output of metrics**

```
sdk_balance_monitor_address_balance{address="0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",name="starknet_strk_user1_account",network="starknet_strk"} 0.029946571772853176
sdk_balance_monitor_address_balance{address="0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",name="starknet_user1_account",network="starknet"} 0.002010528992142664
sdk_balance_monitor_address_balance{address="0x282976F47d7C1cE8B9f956241105d5741f203fA6",name="metis_user1_account",network="evm"} 0.3353651618530689
sdk_balance_monitor_address_balance{address="0x81251B395506b4BfF35DB657751285c138BFB3F5",name="metis_user2_account",network="evm"} 0.23918426945428795

sdk_balance_monitor_balance_empty{address="0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",name="starknet_strk_user1_account",network="starknet_strk"} 0
sdk_balance_monitor_balance_empty{address="0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",name="starknet_user1_account",network="starknet"} 0
sdk_balance_monitor_balance_empty{address="0x282976F47d7C1cE8B9f956241105d5741f203fA6",name="metis_user1_account",network="evm"} 0
sdk_balance_monitor_balance_empty{address="0x81251B395506b4BfF35DB657751285c138BFB3F5",name="metis_user2_account",network="evm"} 0

sdk_balance_monitor_balance_low{address="0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",name="starknet_strk_user1_account",network="starknet_strk"} 0
sdk_balance_monitor_balance_low{address="0x07892c7bf9fc9b9135ab18ddfdc3640d75aa56ed9761fa43bf7eda8bfdfc5919",name="starknet_user1_account",network="starknet"} 0
sdk_balance_monitor_balance_low{address="0x282976F47d7C1cE8B9f956241105d5741f203fA6",name="metis_user1_account",network="evm"} 0
sdk_balance_monitor_balance_low{address="0x81251B395506b4BfF35DB657751285c138BFB3F5",name="metis_user2_account",network="evm"} 0

sdk_balance_monitor_rpc_bad{name="metis_user1_account",network="evm"} 0
sdk_balance_monitor_rpc_bad{name="metis_user2_account",network="evm"} 0
sdk_balance_monitor_rpc_bad{name="starknet_strk_user1_account",network="starknet_strk"} 0
sdk_balance_monitor_rpc_bad{name="starknet_user1_account",network="starknet"} 0
```