# evmlab

A local Ethereum development environment with a Go HTTP API for low-level EVM operations — transaction building, signing, broadcasting, RPC queries, and utility tooling. Includes a Solidity contract suite for experimentation.

## Stack

- **Local chain** — private geth node + Blockscout explorer via Docker Compose
- **Go API** (`ipx/`) — chi-based HTTP server with Swagger UI
- **Solidity** (`contracts/`) — sample contracts (currently: multi-account vault)
- **Rust EVM Engine** *(planned)* — ABI encoding/decoding, calldata construction, EVM simulation, gas estimation, revert decoding, event decoding, policy warnings

## API

| Group | Base path   | Description                                                                                                   |
|-------|-------------|---------------------------------------------------------------------------------------------------------------|
| RPC   | `/evm/rpc`  | Chain queries: chain-id, block-number, gas-price, nonce, balance, estimate-gas, eth-call, transaction/receipt |
| v1    | `/evm/v1`   | Hash (keccak256), sign/verify, legacy transaction build & sign                                                |
| v2    | `/evm/v2`   | EIP-1559 native transfer                                                                                      |
| Tool  | `/evm/tool` | EIP-55 checksum, key derivation from private key                                                              |

Swagger UI available at `http://localhost:33152/swagger/index.html`.

## Quick start

```bash
# Start local geth + Blockscout
make up

# Build and run the API server
make server

# Deploy a contract
make deploy CONTRACT=contracts/vault/MultiAccountVault.sol DEPLOYER=0xYourAddress
```

## Configuration

Copy `config.example.yaml` to `config.yaml` and adjust as needed.

> **Warning**: The example config contains pre-funded local dev keys. Never use them on any public network.
