package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultRPCURL = "http://127.0.0.1:8545"
	masterKeyHex  = "ea66255f7b410dd47cbe9c9c9bf049bfe441877c1eaf09c5faef8ab7cf45357d"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	rpcURL := getenv("RPC_URL", defaultRPCURL)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("connect rpc: %w", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get chain id: %w", err)
	}

	privateKey, from, err := privateKeyToAddress(masterKeyHex)
	if err != nil {
		return err
	}
	fmt.Println("Deployer:", from.Hex())

	bytecode, err := loadBytecode()
	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return fmt.Errorf("get nonce: %w", err)
	}

	fmt.Println("Estimate From:", from.Hex())
	fmt.Println("Bytecode size:", len(bytecode))

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    nil,
		Value: big.NewInt(0),
		Data:  bytecode,
	})
	if err != nil {
		return fmt.Errorf("estimate gas: %w", err)
	}
	gasLimit = gasLimit + gasLimit/5
	tipCap, err := client.SuggestGasTipCap(ctx)

	if err != nil {
		return fmt.Errorf("suggest gas tip cap: %w", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)

	if err != nil {
		return fmt.Errorf("get latest header: %w", err)
	}

	feeCap := new(big.Int).Add(
		new(big.Int).Mul(header.BaseFee, big.NewInt(2)),
		tipCap,
	)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        nil,
		Value:     big.NewInt(0),
		Data:      bytecode,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)

	if err != nil {
		return fmt.Errorf("sign tx: %w", err)
	}

	if err = client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("send tx: %w", err)
	}

	fmt.Println("RPC:", rpcURL)
	fmt.Println("Chain ID:", chainID.String())
	fmt.Println("Deployer:", from.Hex())
	fmt.Println("Deploy tx:", signedTx.Hash().Hex())

	receipt, err := waitReceipt(ctx, client, signedTx.Hash(), 30*time.Second)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("deployment failed: status=%d", receipt.Status)
	}

	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	if err != nil {
		return fmt.Errorf("get deployed code: %w", err)
	}

	fmt.Println("Contract address:", receipt.ContractAddress.Hex())
	fmt.Println("Block number:", receipt.BlockNumber.String())
	fmt.Println("Gas used:", receipt.GasUsed)
	fmt.Println("Runtime code size:", len(code), "bytes")

	return nil
}

func loadBytecode() ([]byte, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	binPath := filepath.Join(root, "build", "contracts_MultiAccountVault_sol_MultiAccountVault.bin")
	binBytes, err := os.ReadFile(binPath)
	if err != nil {
		return nil, fmt.Errorf("read bytecode: %w", err)
	}

	binHex := strings.TrimSpace(string(binBytes))
	binHex = strings.TrimPrefix(binHex, "0x")

	bytecode, err := hex.DecodeString(binHex)
	if err != nil {
		return nil, fmt.Errorf("decode bytecode: %w", err)
	}

	return bytecode, nil
}

func privateKeyToAddress(keyHex string) (*ecdsa.PrivateKey, common.Address, error) {
	keyHex = strings.TrimPrefix(strings.TrimSpace(keyHex), "0x")

	privateKey, err := gethcrypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("parse private key: %w", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("invalid public key")
	}

	return privateKey, gethcrypto.PubkeyToAddress(*publicKey), nil
}

func waitReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	deadline := time.Now().Add(timeout)

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for receipt: %s", txHash.Hex())
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if exists(filepath.Join(dir, "build")) && exists(filepath.Join(dir, "contracts")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found from %s", wd)
		}

		dir = parent
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
