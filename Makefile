PROJECT = evmlab
DOCKER ?= docker

GETH_CONTAINER = evmlab-geth
GETH_VOLUME = evmlab-geth-data
BLOCKSCOUT_DB_VOLUME = evmlab-blockscout-db-data

BIN_DIR = bin
BUILD_DIR = build

GO_DIR = ipx
GO_DEPLOYER_BIN = contract_deployer
GO_SERVER_BIN   = server

.PHONY: help \
	up down logs reset \
	geth-logs geth-attach geth-reset \
	compile deploy test \
	deploy-lab-token deploy-vault \
	go-build-deployer go-build-server go-build \
	server go-test swag \
	clean

help:
	@echo "$(PROJECT) commands:"
	@echo ""
	@echo "  make up              Start local geth and Blockscout"
	@echo "  make down            Stop local stack"
	@echo "  make logs            Follow docker compose logs"
	@echo "  make reset           Remove local geth and Blockscout data, then restart"
	@echo ""
	@echo "  make geth-logs       Follow geth logs"
	@echo "  make geth-attach     Attach to geth IPC console"
	@echo "  make geth-reset      Remove local geth volume and restart"
	@echo ""
	@echo "  make compile              Compile Solidity contract (CONTRACT=<path>)"
	@echo "  make deploy               Compile and deploy to local geth (CONTRACT=<path>)"
	@echo "  make deploy-lab-token     Deploy LabToken (DEPLOYER=key NAME=... SYMBOL=... DECIMALS=... SUPPLY=...)"
	@echo "  make deploy-vault         Deploy MultiAccountVault (DEPLOYER=key)"
	@echo ""
	@echo "  make go-build          Build all Go binaries"
	@echo "  make go-build-deployer Build contract_deployer binary"
	@echo "  make go-build-server   Build server binary"
	@echo "  make server            Build and run API server"
	@echo "  make go-test           Run Go tests"
	@echo "  make swag            Regenerate Swagger docs"
	@echo ""
	@echo "  make clean           Remove generated build outputs"

up:
	$(DOCKER) compose up -d --wait

down:
	$(DOCKER) compose down

logs:
	$(DOCKER) compose logs -f

reset:
	$(DOCKER) compose down -v
	$(DOCKER) compose up -d --wait

geth-logs:
	$(DOCKER) logs -f $(GETH_CONTAINER)

geth-attach:
	$(DOCKER) exec -it $(GETH_CONTAINER) geth attach /root/.ethereum/geth.ipc

geth-reset:
	$(DOCKER) compose down
	$(DOCKER) volume rm $(GETH_VOLUME) || true
	$(DOCKER) compose up -d --wait

CONTRACT_DIR    = $(shell dirname $(CONTRACT))
CONTRACT_NAME   = $(shell basename $(CONTRACT) .sol)
CONTRACT_SUBDIR = $(CONTRACT_NAME)

compile:
	@[ -n "$(CONTRACT)" ] || (echo "error: CONTRACT is required (e.g. make compile CONTRACT=contracts/vault/MultiAccountVault.sol)" && exit 1)
	mkdir -p $(BUILD_DIR)/$(CONTRACT_SUBDIR)
	npm run compile -- -o $(BUILD_DIR)/$(CONTRACT_SUBDIR) $(CONTRACT)
	npm run standard-json -- $(CONTRACT)

deploy: go-build-deployer compile
	@[ -n "$(DEPLOYER)" ] || (echo "error: DEPLOYER is required (e.g. make deploy CONTRACT=... DEPLOYER=key0)" && exit 1)
	./$(BIN_DIR)/$(GO_DEPLOYER_BIN) --contract $(CONTRACT) --deployer $(DEPLOYER)

LAB_TOKEN_CONTRACT = contracts/lab_token/LabToken.sol

# e.g. make deploy-lab-token DEPLOYER=0xEbD69375d51a8472DF22A3C18405b5A2586c2Aa2 NAME=LabToken SYMBOL=LAB DECIMALS=6 SUPPLY=10000000000
deploy-lab-token: go-build-deployer
	@[ -n "$(DEPLOYER)" ] || (echo "error: DEPLOYER is required" && exit 1)
	@[ -n "$(NAME)" ]     || (echo "error: NAME is required"     && exit 1)
	@[ -n "$(SYMBOL)" ]   || (echo "error: SYMBOL is required"   && exit 1)
	@[ -n "$(DECIMALS)" ] || (echo "error: DECIMALS is required" && exit 1)
	@[ -n "$(SUPPLY)" ]   || (echo "error: SUPPLY is required"   && exit 1)
	$(MAKE) compile CONTRACT=$(LAB_TOKEN_CONTRACT)
	./$(BIN_DIR)/$(GO_DEPLOYER_BIN) \
		--contract $(LAB_TOKEN_CONTRACT) \
		--deployer $(DEPLOYER) \
		--ctor '{"types":["string","string","uint8","uint256"],"args":["$(NAME)","$(SYMBOL)","$(DECIMALS)","$(SUPPLY)"]}'

VAULT_CONTRACT = contracts/vault/MultiAccountVault.sol

# e.g. make deploy-vault DEPLOYER=0xEbD69375d51a8472DF22A3C18405b5A2586c2Aa2
deploy-vault: go-build-deployer
	@[ -n "$(DEPLOYER)" ] || (echo "error: DEPLOYER is required" && exit 1)
	$(MAKE) compile CONTRACT=$(VAULT_CONTRACT)
	./$(BIN_DIR)/$(GO_DEPLOYER_BIN) \
		--contract $(VAULT_CONTRACT) \
		--deployer $(DEPLOYER)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

go-build-deployer: $(BIN_DIR)
	cd $(GO_DIR) && go build -o ../$(BIN_DIR)/$(GO_DEPLOYER_BIN) ./cmd/contract_deployer

go-build-server: $(BIN_DIR)
	cd $(GO_DIR) && go build -o ../$(BIN_DIR)/$(GO_SERVER_BIN) ./cmd/server

go-build: go-build-deployer go-build-server

server: go-build-server swag
	./$(BIN_DIR)/$(GO_SERVER_BIN)

go-test:
	cd $(GO_DIR) && go test ./...

swag:
	cd $(GO_DIR) && swag init -g cmd/server/main.go -o docs

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	rm -rf artifacts cache
