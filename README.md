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

| Method | Path                           | Description                     |
|--------|--------------------------------|---------------------------------|
| POST   | `/evm/rpc/`                    | Raw JSON-RPC proxy              |
| POST   | `/evm/rpc/chain-id`            | eth_chainId                     |
| POST   | `/evm/rpc/block-number`        | eth_blockNumber                 |
| POST   | `/evm/rpc/nonce`               | eth_getTransactionCount         |
| POST   | `/evm/rpc/balance`             | eth_getBalance                  |
| POST   | `/evm/rpc/code`                | eth_getCode                     |
| POST   | `/evm/rpc/transaction`         | eth_getTransactionByHash        |
| POST   | `/evm/rpc/transaction/receipt` | eth_getTransactionReceipt       |
| POST   | `/evm/rpc/transaction/send`    | eth_sendRawTransaction          |
| POST   | `/evm/rpc/transaction/status`  | Transaction status              |
| POST   | `/evm/rpc/batch`               | Batch JSON-RPC proxy            |
| POST   | `/evm/rpc/fee/base`            | baseFeePerGas from latest block |
| POST   | `/evm/rpc/fee/priority`        | eth_maxPriorityFeePerGas        |
| POST   | `/evm/rpc/fee/max`             | maxFeePerGas (2×base + tip)     |
| POST   | `/evm/rpc/gas/price`           | eth_gasPrice                    |
| POST   | `/evm/rpc/gas/estimate`        | eth_estimateGas                 |
| POST   | `/evm/rpc/call`                | eth_call                        |

### ABI — `/evm/abi`

| Method | Path                            | Description                                      |
|--------|---------------------------------|--------------------------------------------------|
| POST   | `/evm/abi/selector`             | 4-byte selector from function signature          |
| POST   | `/evm/abi/encode`               | ABI-encode calldata from signature + args        |
| POST   | `/evm/abi/decode/result`        | Decode eth_call return data by ABI types         |
| POST   | `/evm/abi/decode/call`          | Decode calldata by function signature            |
| POST   | `/evm/abi/decode/revert`        | Decode revert data by error signature            |
| POST   | `/evm/abi/encode/eip712-domain` | `EIP712Domain` calldata                          |

### Contract — `/evm/contract`

| Method | Path                                         | Description                                      |
|--------|----------------------------------------------|--------------------------------------------------|
| POST   | `/evm/contract/multicall3/aggregate3`        | Multicall3 aggregate3 batch eth_call             |
| POST   | `/evm/contract/eip712/domain`                | Fetch EIP-712 domain from contract               |
| POST   | `/evm/contract/eip2612/nonces`               | Fetch EIP-2612 permit nonce                      |
| POST   | `/evm/contract/erc20/detect`                 | Heuristic ERC-20-like detection                  |
| POST   | `/evm/contract/erc20/metadata`               | Fetch ERC-20 name, symbol, decimals, totalSupply |
| POST   | `/evm/contract/erc20/balance`                | `balanceOf(address)`                             |
| POST   | `/evm/contract/erc20/allowance`              | `allowance(address,address)`                     |
| POST   | `/evm/contract/erc20/approved`               | Check if allowance >= amount                     |
| POST   | `/evm/contract/erc20/calldata/balance-of`    | `balanceOf(address)` calldata                    |
| POST   | `/evm/contract/erc20/calldata/approve`       | `approve(address,uint256)` calldata              |
| POST   | `/evm/contract/erc20/calldata/transfer`      | `transfer(address,uint256)` calldata             |
| POST   | `/evm/contract/erc20/calldata/allowance`     | `allowance(address,address)` calldata            |
| POST   | `/evm/contract/erc20/calldata/transfer-from` | `transferFrom(address,address,uint256)` calldata |

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

| Method | Path                                        | Description                                       |
|--------|---------------------------------------------|---------------------------------------------------|
| POST   | `/evm/v2/transaction/native/legacy`         | Unsigned legacy native transfer                   |
| POST   | `/evm/v2/transaction/native/eip1559`        | Unsigned EIP-1559 native transfer                 |
| POST   | `/evm/v2/transaction/erc20/legacy`          | Unsigned legacy ERC-20 transfer (gas estimated)   |
| POST   | `/evm/v2/transaction/erc20/eip1559`         | Unsigned EIP-1559 ERC-20 transfer (gas estimated) |
| POST   | `/evm/v2/transaction/contract/legacy`       | Unsigned legacy contract call                     |
| POST   | `/evm/v2/transaction/contract/eip1559`      | Unsigned EIP-1559 contract call                   |

### v3 — `/evm/v3` — build and sign tx

| Method | Path                                        | Description                    |
|--------|---------------------------------------------|--------------------------------|
| POST   | `/evm/v3/transaction/native/legacy`         | Sign legacy native transfer    |
| POST   | `/evm/v3/transaction/native/eip1559`        | Sign EIP-1559 native transfer  |
| POST   | `/evm/v3/transaction/erc20/legacy`          | Sign legacy ERC-20 transfer    |
| POST   | `/evm/v3/transaction/erc20/eip1559`         | Sign EIP-1559 ERC-20 transfer  |
| POST   | `/evm/v3/transaction/contract/legacy`       | Sign legacy contract call      |
| POST   | `/evm/v3/transaction/contract/eip1559`      | Sign EIP-1559 contract call    |

### v4 — `/evm/v4` — build, sign, and broadcast tx

| Method | Path                                        | Description                                 |
|--------|---------------------------------------------|---------------------------------------------|
| POST   | `/evm/v4/transaction/native/legacy`         | Sign and broadcast legacy native transfer   |
| POST   | `/evm/v4/transaction/native/eip1559`        | Sign and broadcast EIP-1559 native transfer |
| POST   | `/evm/v4/transaction/erc20/legacy`          | Sign and broadcast legacy ERC-20 transfer   |
| POST   | `/evm/v4/transaction/erc20/eip1559`         | Sign and broadcast EIP-1559 ERC-20 transfer |
| POST   | `/evm/v4/transaction/contract/legacy`       | Sign and broadcast legacy contract call     |
| POST   | `/evm/v4/transaction/contract/eip1559`      | Sign and broadcast EIP-1559 contract call   |

### Sanctum — `/evm/sanctum` — build, sign, and broadcast Sanctum contract calls

#### Nexus — `/evm/sanctum/nexus` — membership management

| Method | Path                                              | Description                                        |
|--------|---------------------------------------------------|----------------------------------------------------|
| POST   | `/evm/sanctum/nexus/register/legacy`              | `register()` — legacy                              |
| POST   | `/evm/sanctum/nexus/register/eip1559`             | `register()` — EIP-1559                            |
| POST   | `/evm/sanctum/nexus/register/for/legacy`          | `registerFor(address)` — legacy                    |
| POST   | `/evm/sanctum/nexus/register/for/eip1559`         | `registerFor(address)` — EIP-1559                  |
| POST   | `/evm/sanctum/nexus/register/approve/legacy`      | `approveRegister(address)` — legacy                |
| POST   | `/evm/sanctum/nexus/register/approve/eip1559`     | `approveRegister(address)` — EIP-1559              |
| POST   | `/evm/sanctum/nexus/deregister/legacy`            | `deregister()` — legacy                            |
| POST   | `/evm/sanctum/nexus/deregister/eip1559`           | `deregister()` — EIP-1559                          |
| POST   | `/evm/sanctum/nexus/deregister/for/legacy`        | `deregisterFor(address)` — legacy                  |
| POST   | `/evm/sanctum/nexus/deregister/for/eip1559`       | `deregisterFor(address)` — EIP-1559                |
| POST   | `/evm/sanctum/nexus/account/list`                 | `getAccounts()` — list all registered accounts     |
| POST   | `/evm/sanctum/nexus/account/count`                | `accountCount()` — number of registered accounts   |
| POST   | `/evm/sanctum/nexus/account/info`                 | `getAccountInfo(address)` — role and block info    |

#### Treasury — `/evm/sanctum/treasury` — native token treasury operations

| Method | Path                                                | Description                                            |
|--------|-----------------------------------------------------|--------------------------------------------------------|
| POST   | `/evm/sanctum/treasury/native/deposit/legacy`       | `depositNative()` payable — legacy                     |
| POST   | `/evm/sanctum/treasury/native/deposit/eip1559`      | `depositNative()` payable — EIP-1559                   |
| POST   | `/evm/sanctum/treasury/native/request/legacy`       | `requestNative(uint256)` — legacy                      |
| POST   | `/evm/sanctum/treasury/native/request/eip1559`      | `requestNative(uint256)` — EIP-1559                    |
| POST   | `/evm/sanctum/treasury/native/approve/legacy`       | `approveNative(address,uint256)` — legacy              |
| POST   | `/evm/sanctum/treasury/native/approve/eip1559`      | `approveNative(address,uint256)` — EIP-1559            |
| POST   | `/evm/sanctum/treasury/native/approve/all/legacy`   | `approveNativeAll(address)` — legacy                   |
| POST   | `/evm/sanctum/treasury/native/approve/all/eip1559`  | `approveNativeAll(address)` — EIP-1559                 |
| POST   | `/evm/sanctum/treasury/native/withdraw/legacy`      | `withdrawNative(uint256)` — legacy                     |
| POST   | `/evm/sanctum/treasury/native/withdraw/eip1559`     | `withdrawNative(uint256)` — EIP-1559                   |
| POST   | `/evm/sanctum/treasury/native/withdraw/all/legacy`  | `withdrawNativeAll()` — legacy                         |
| POST   | `/evm/sanctum/treasury/native/withdraw/all/eip1559` | `withdrawNativeAll()` — EIP-1559                       |
| POST   | `/evm/sanctum/treasury/native/balance`              | `nativeBalance()` — total treasury balance             |
| POST   | `/evm/sanctum/treasury/native/available`            | `nativeAvailable()` — unallocated balance              |
| POST   | `/evm/sanctum/treasury/native/allocation`           | `nativeAllocation(address)` — approved amount for user |
| POST   | `/evm/sanctum/treasury/native/pending`              | `nativePending(address)` — pending request for user    |

---

## Configuration

Copy `config.example.yaml` to `config.yaml` and adjust as needed.

> **Warning**: The example config contains pre-funded local dev keys. Never use them on any public network.