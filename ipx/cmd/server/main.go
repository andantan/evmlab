// @title           evmlab API
// @version         1.0
// @description     Ethereum transaction API
// @host            localhost:33152
// @BasePath        /
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andantan/evmlab/api/handler/contract"
	"github.com/andantan/evmlab/api/handler/misc"
	"github.com/andantan/evmlab/api/handler/sanctum"
	"github.com/andantan/evmlab/api/handler/v1"
	"github.com/andantan/evmlab/api/handler/v2"
	"github.com/andantan/evmlab/api/handler/v3"
	"github.com/andantan/evmlab/api/handler/v4"
	_ "github.com/andantan/evmlab/docs"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := util.FindProjectRoot()
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(root, "config.yaml"))
	if err != nil {
		return err
	}

	client := rpc.NewClient(cfg.RPCURL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/evm/rpc", func(r chi.Router) {
		rpcHandler := misc.NewRPCHandler(client)
		r.Post("/", rpcHandler.Raw)
		r.Post("/chain-id", rpcHandler.ChainID)
		r.Post("/block-number", rpcHandler.BlockNumber)
		r.Post("/nonce", rpcHandler.Nonce)
		r.Post("/balance", rpcHandler.Balance)
		r.Post("/code", rpcHandler.Code)
		r.Post("/transaction", rpcHandler.Transaction)
		r.Post("/transaction/receipt", rpcHandler.TransactionReceipt)
		r.Post("/transaction/send", rpcHandler.SendTransaction)
		r.Post("/transaction/status", rpcHandler.TransactionStatus)
		r.Post("/batch", rpcHandler.Batch)
		r.Post("/fee/base", rpcHandler.BaseFeePerGas)
		r.Post("/fee/priority", rpcHandler.MaxPriorityFeePerGas)
		r.Post("/fee/max", rpcHandler.MaxFeePerGas)
		r.Post("/gas/price", rpcHandler.GasPrice)
		r.Post("/gas/estimate", rpcHandler.EstimateGas)
		r.Post("/call", rpcHandler.Call)
	})

	r.Route("/evm/abi", func(r chi.Router) {
		abi := misc.NewAbiHandler()
		r.Post("/selector", abi.Selector)
		r.Post("/decode/result", abi.DecodeResult)
		r.Post("/decode/call", abi.DecodeCall)
		r.Post("/decode/revert", abi.DecodeRevert)
		r.Post("/encode", abi.Encode)
		r.Post("/encode/eip712-domain", abi.EIP712DomainCalldata)
	})

	r.Route("/evm/tool", func(r chi.Router) {
		tool := misc.NewToolHandler()
		r.Post("/address/eip55", tool.EIP55)
		r.Post("/crypto/derive", tool.DeriveKey)
		r.Post("/unit/convert", tool.ConvertUnit)
	})

	r.Route("/evm/hash", func(r chi.Router) {
		hash := misc.NewHashHandler()
		r.Post("/keccak256/legacy", hash.Keccak256Legacy)
		r.Post("/keccak256/eip191", hash.Keccak256EIP191)
		r.Post("/keccak256/eip712", hash.Keccak256EIP712)
	})

	r.Route("/evm/sign", func(r chi.Router) {
		sign := misc.NewSignHandler(cfg)
		r.Post("/", sign.Sign)
		r.Post("/ecrecover", sign.Ecrecover)
		r.Post("/verify/by-public-key", sign.VerifyByPublicKey)
		r.Post("/verify/by-address", sign.VerifyByAddress)
		r.Post("/transaction/legacy", sign.SignLegacyTransaction)
		r.Post("/transaction/eip1559", sign.SignEIP1559Transaction)
	})

	r.Route("/evm/contract", func(r chi.Router) {
		mc3 := contract.NewMulticall3Handler(cfg, client)
		r.Post("/multicall3/aggregate3", mc3.Aggregate3)

		eip := contract.NewEIPHandler(client)
		r.Post("/eip712/domain", eip.EIP712Domain)
		r.Post("/eip2612/nonces", eip.EIP2612Nonces)

		erc20 := contract.NewERC20Handler(cfg, client)
		r.Post("/erc20/detect", erc20.Detect)
		r.Post("/erc20/metadata", erc20.Metadata)
		r.Post("/erc20/balance", erc20.Balance)
		r.Post("/erc20/allowance", erc20.Allowance)
		r.Post("/erc20/approved", erc20.Approved)
		r.Post("/erc20/calldata/balance-of", erc20.BalanceOfCalldata)
		r.Post("/erc20/calldata/approve", erc20.ApproveCalldata)
		r.Post("/erc20/calldata/transfer", erc20.TransferCalldata)
		r.Post("/erc20/calldata/allowance", erc20.AllowanceCalldata)
		r.Post("/erc20/calldata/transfer-from", erc20.TransferFromCalldata)
	})

	r.Route("/evm/v1", func(r chi.Router) {
		tx := v1.NewTransactionHandler(cfg)
		r.Post("/transaction/legacy/build", tx.BuildLegacyTransaction)
		r.Post("/transaction/eip1559/build", tx.BuildEIP1559Transaction)
	})

	r.Route("/evm/v2", func(r chi.Router) {
		transfer := v2.NewTransactionHandler(client)
		r.Post("/transaction/native/legacy", transfer.BuildNativeLegacyTransaction)
		r.Post("/transaction/native/eip1559", transfer.BuildNativeEIP1559Transaction)
		r.Post("/transaction/erc20/legacy", transfer.BuildERC20LegacyTransaction)
		r.Post("/transaction/erc20/eip1559", transfer.BuildERC20EIP1559Transaction)

		r.Post("/transaction/contract/legacy", transfer.BuildContractCallLegacyTransaction)
		r.Post("/transaction/contract/eip1559", transfer.BuildContractCallEIP1559Transaction)
	})

	r.Route("/evm/v3", func(r chi.Router) {
		tx := v3.NewTransactionHandler(cfg, client)
		r.Post("/transaction/native/legacy", tx.BuildNativeLegacyTransaction)
		r.Post("/transaction/native/eip1559", tx.BuildNativeEIP1559Transaction)
		r.Post("/transaction/erc20/legacy", tx.BuildERC20LegacyTransaction)
		r.Post("/transaction/erc20/eip1559", tx.BuildERC20EIP1559Transaction)
		r.Post("/transaction/contract/legacy", tx.BuildContractCallLegacyTransaction)
		r.Post("/transaction/contract/eip1559", tx.BuildContractCallEIP1559Transaction)
	})

	r.Route("/evm/v4", func(r chi.Router) {
		tx := v4.NewTransactionHandler(cfg, client)
		r.Post("/transaction/native/legacy", tx.BuildNativeLegacyTransaction)
		r.Post("/transaction/native/eip1559", tx.BuildNativeEIP1559Transaction)
		r.Post("/transaction/erc20/legacy", tx.BuildERC20LegacyTransaction)
		r.Post("/transaction/erc20/eip1559", tx.BuildERC20EIP1559Transaction)
		r.Post("/transaction/contract/legacy", tx.BuildContractCallLegacyTransaction)
		r.Post("/transaction/contract/eip1559", tx.BuildContractCallEIP1559Transaction)
	})

	r.Route("/evm/sanctum", func(r chi.Router) {
		s := sanctum.NewSanctumHandler(cfg, client)
		r.Route("/nexus", func(r chi.Router) {
			r.Post("/register/legacy", s.RegisterLegacy)
			r.Post("/register/eip1559", s.RegisterEIP1559)
			r.Post("/register/for/legacy", s.RegisterForLegacy)
			r.Post("/register/for/eip1559", s.RegisterForEIP1559)
			r.Post("/register/approve/legacy", s.ApproveRegisterLegacy)
			r.Post("/register/approve/eip1559", s.ApproveRegisterEIP1559)
			r.Post("/deregister/legacy", s.DeregisterLegacy)
			r.Post("/deregister/eip1559", s.DeregisterEIP1559)
			r.Post("/deregister/for/legacy", s.DeregisterForLegacy)
			r.Post("/deregister/for/eip1559", s.DeregisterForEIP1559)

			r.Post("/account/list", s.GetAccounts)
			r.Post("/account/count", s.AccountCount)
			r.Post("/account/info", s.GetAccountInfo)
		})
		r.Route("/treasury", func(r chi.Router) {
			r.Post("/native/deposit/legacy", s.DepositNativeLegacy)
			r.Post("/native/deposit/eip1559", s.DepositNativeEIP1559)
			r.Post("/native/request/legacy", s.RequestNativeLegacy)
			r.Post("/native/request/eip1559", s.RequestNativeEIP1559)
			r.Post("/native/approve/legacy", s.ApproveNativeLegacy)
			r.Post("/native/approve/eip1559", s.ApproveNativeEIP1559)
			r.Post("/native/approve/all/legacy", s.ApproveNativeAllLegacy)
			r.Post("/native/approve/all/eip1559", s.ApproveNativeAllEIP1559)
			r.Post("/native/withdraw/legacy", s.WithdrawNativeLegacy)
			r.Post("/native/withdraw/eip1559", s.WithdrawNativeEIP1559)
			r.Post("/native/withdraw/all/legacy", s.WithdrawNativeAllLegacy)
			r.Post("/native/withdraw/all/eip1559", s.WithdrawNativeAllEIP1559)
			r.Post("/native/balance", s.NativeBalance)
			r.Post("/native/available", s.NativeAvailable)
			r.Post("/native/allocation", s.NativeAllocation)
			r.Post("/native/pending", s.NativePending)
		})
	})

	fmt.Println("Listening on", cfg.ServerAddr)
	return http.ListenAndServe(cfg.ServerAddr, r)
}
