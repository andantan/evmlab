#!/usr/bin/env bash
set -euo pipefail

CONTRACT_NAME="${1:?usage: $0 <ContractName>}"
export CONTRACT_NAME
mkdir -p build

node <<'NODE'
const fs = require("fs");

const files = [
  "contracts/vault/MultiAccountVault.sol",
  "contracts/vault/abstract/VaultAccess.sol",
  "contracts/vault/abstract/VaultAccounts.sol",
  "contracts/vault/abstract/VaultFunds.sol",
  "contracts/vault/abstract/VaultDistribution.sol",
  "contracts/vault/interfaces/IMultiAccountVault.sol",
  "contracts/vault/interfaces/IVaultAccounts.sol",
  "contracts/vault/interfaces/IVaultFunds.sol",
  "contracts/vault/interfaces/IVaultDistribution.sol",
  "contracts/vault/libraries/VaultTypes.sol"
];

const sources = {};

for (const file of files) {
  sources[file] = {
    content: fs.readFileSync(file, "utf8")
  };
}

const input = {
  language: "Solidity",
  sources,
  settings: {
    optimizer: {
      enabled: false,
      runs: 200
    },
    outputSelection: {
      "*": {
        "*": [
          "abi",
          "evm.bytecode",
          "evm.deployedBytecode",
          "metadata"
        ]
      }
    }
  }
};

fs.writeFileSync(`build/vault/${process.env.CONTRACT_NAME}.standard-input.json`, JSON.stringify(input, null, 2));
NODE

echo "wrote build/vault/${CONTRACT_NAME}.standard-input.json"