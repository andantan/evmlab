# evmlab

A local Ethereum development environment for low-level EVM operations.

- **Go** (`ipx/`) — HTTP API server: transaction building, signing, broadcasting, ABI encoding, RPC queries, and utility tooling
- **Rust** *(planned)* — CLI: ABI encoding/decoding, calldata construction, EVM simulation, gas estimation, revert/event decoding

## Stack

- **Local chain** — private geth node + Blockscout explorer via Docker Compose
- **Go API** (`ipx/`) — chi-based HTTP server with Swagger UI
- **Solidity** (`contracts/`) — sample contracts (vault, ERC-20 token)

---

## Make commands

### Docker / local chain

| Command             | Description                                   |
|---------------------|-----------------------------------------------|
| `make up`           | Start local geth and Blockscout               |
| `make down`         | Stop local stack                              |
| `make logs`         | Follow docker compose logs                    |
| `make reset`        | Remove all local data volumes and restart     |
| `make geth-logs`    | Follow geth container logs                    |
| `make geth-attach`  | Attach to geth IPC console                    |
| `make geth-reset`   | Remove geth volume only and restart           |

### Contracts

| Command                                                                            | Description                                              |
|------------------------------------------------------------------------------------|----------------------------------------------------------|
| `make compile CONTRACT=<path>`                                                     | Compile a Solidity file and generate standard JSON input |
| `make deploy CONTRACT=<path> DEPLOYER=<key>`                                       | Compile and deploy to local geth                         |
| `make deploy-lab-token DEPLOYER=<key> NAME=<n> SYMBOL=<s> DECIMALS=<d> SUPPLY=<n>` | Deploy LabToken ERC-20                                   |
| `make deploy-vault DEPLOYER=<key>`                                                 | Deploy MultiAccountVault                                 |

### Go

| Command                  | Description                            |
|--------------------------|----------------------------------------|
| `make server`            | Build and run the API server           |
| `make go-build`          | Build all Go binaries                  |
| `make go-build-server`   | Build server binary only               |
| `make go-build-deployer` | Build contract_deployer binary only    |
| `make go-test`           | Run Go tests                           |
| `make swag`              | Regenerate Swagger docs                |
| `make clean`             | Remove `bin/` and `build/` directories |

---

## API endpoints

Swagger UI: `http://localhost:33152/swagger/index.html`

### RPC — `/evm/rpc`

| Method | Path                           | Description                      |
|--------|--------------------------------|----------------------------------|
| POST   | `/evm/rpc/`                    | Raw JSON-RPC proxy               |
| POST   | `/evm/rpc/chain-id`            | eth_chainId                      |
| POST   | `/evm/rpc/block-number`        | eth_blockNumber                  |
| POST   | `/evm/rpc/nonce`               | eth_getTransactionCount          |
| POST   | `/evm/rpc/balance`             | eth_getBalance                   |
| POST   | `/evm/rpc/code`                | eth_getCode                      |
| POST   | `/evm/rpc/transaction`         | eth_getTransactionByHash         |
| POST   | `/evm/rpc/transaction/receipt` | eth_getTransactionReceipt        |
| POST   | `/evm/rpc/transaction/send`    | eth_sendRawTransaction           |
| POST   | `/evm/rpc/fee/base`            | baseFeePerGas from latest block  |
| POST   | `/evm/rpc/fee/priority`        | eth_maxPriorityFeePerGas         |
| POST   | `/evm/rpc/fee/max`             | maxFeePerGas (2×base + tip)      |
| POST   | `/evm/rpc/gas/price`           | eth_gasPrice                     |
| POST   | `/evm/rpc/gas/estimate`        | eth_estimateGas                  |
| POST   | `/evm/rpc/call`                | eth_call                         |

### ABI — `/evm/abi`

| Method | Path                            | Description                                      |
|--------|---------------------------------|--------------------------------------------------|
| POST   | `/evm/abi/selector`             | 4-byte selector from function signature          |
| POST   | `/evm/abi/encode`               | ABI-encode calldata from signature + args        |
| POST   | `/evm/abi/decode/result`        | Decode eth_call return data by ABI types         |
| POST   | `/evm/abi/decode/call`          | Decode calldata by function signature            |
| POST   | `/evm/abi/encode/balance-of`    | `balanceOf(address)` calldata                    |
| POST   | `/evm/abi/encode/approve`       | `approve(address,uint256)` calldata              |
| POST   | `/evm/abi/encode/transfer`      | `transfer(address,uint256)` calldata             |
| POST   | `/evm/abi/encode/allowance`     | `allowance(address,address)` calldata            |
| POST   | `/evm/abi/encode/transfer-from` | `transferFrom(address,address,uint256)` calldata |

### Hash — `/evm/hash`

| Method | Path                          | Description                    |
|--------|-------------------------------|--------------------------------|
| POST   | `/evm/hash/keccak256/legacy`  | keccak256 of raw data          |
| POST   | `/evm/hash/keccak256/eip191`  | EIP-191 personal_sign hash     |
| POST   | `/evm/hash/keccak256/eip712`  | EIP-712 structured data hash   |

### Sign — `/evm/sign`

| Method | Path                               | Description                           |
|--------|------------------------------------|---------------------------------------|
| POST   | `/evm/sign/`                       | Sign arbitrary hash                   |
| POST   | `/evm/sign/ecrecover`              | Recover signer address from signature |
| POST   | `/evm/sign/verify/by-public-key`   | Verify signature against public key   |
| POST   | `/evm/sign/verify/by-address`      | Verify signature against address      |
| POST   | `/evm/sign/transaction/legacy`     | Sign a pre-built legacy tx            |
| POST   | `/evm/sign/transaction/eip1559`    | Sign a pre-built EIP-1559 tx          |

### Tool — `/evm/tool`

| Method | Path                      | Description                         |
|--------|---------------------------|-------------------------------------|
| POST   | `/evm/tool/address/eip55` | EIP-55 checksum address             |
| POST   | `/evm/tool/crypto/derive` | Derive keypair from private key     |
| POST   | `/evm/tool/unit/convert`  | Wei / Gwei / Ether unit conversion  |

### v1 — `/evm/v1` — build tx from user-supplied fields (no RPC calls)

| Method | Path                                | Description                |
|--------|-------------------------------------|----------------------------|
| POST   | `/evm/v1/transaction/legacy/build`  | Build unsigned legacy tx   |
| POST   | `/evm/v1/transaction/eip1559/build` | Build unsigned EIP-1559 tx |

### v2 — `/evm/v2` — build unsigned tx (fetches chain state, estimates gas)

| Method | Path                                 | Description                                       |
|--------|--------------------------------------|---------------------------------------------------|
| POST   | `/evm/v2/transaction/native/legacy`  | Unsigned legacy native transfer                   |
| POST   | `/evm/v2/transaction/native/eip1559` | Unsigned EIP-1559 native transfer                 |
| POST   | `/evm/v2/transaction/erc20/legacy`   | Unsigned legacy ERC-20 transfer (gas estimated)   |
| POST   | `/evm/v2/transaction/erc20/eip1559`  | Unsigned EIP-1559 ERC-20 transfer (gas estimated) |

### v3 — `/evm/v3` — build and sign tx

| Method | Path                                 | Description                    |
|--------|--------------------------------------|--------------------------------|
| POST   | `/evm/v3/transaction/native/legacy`  | Sign legacy native transfer    |
| POST   | `/evm/v3/transaction/native/eip1559` | Sign EIP-1559 native transfer  |
| POST   | `/evm/v3/transaction/erc20/legacy`   | Sign legacy ERC-20 transfer    |
| POST   | `/evm/v3/transaction/erc20/eip1559`  | Sign EIP-1559 ERC-20 transfer  |

### v4 — `/evm/v4` — build, sign, and broadcast tx

| Method | Path                                 | Description                                 |
|--------|--------------------------------------|---------------------------------------------|
| POST   | `/evm/v4/transaction/native/legacy`  | Sign and broadcast legacy native transfer   |
| POST   | `/evm/v4/transaction/native/eip1559` | Sign and broadcast EIP-1559 native transfer |
| POST   | `/evm/v4/transaction/erc20/legacy`   | Sign and broadcast legacy ERC-20 transfer   |
| POST   | `/evm/v4/transaction/erc20/eip1559`  | Sign and broadcast EIP-1559 ERC-20 transfer |

---

## Configuration

Copy `config.example.yaml` to `config.yaml` and adjust as needed.

> **Warning**: The example config contains pre-funded local dev keys. Never use them on any public network.