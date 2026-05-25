package sanctum

import (
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
)

type SanctumHandler struct {
	cfg    *config.Config
	client *rpc.Client
}

func NewSanctumHandler(cfg *config.Config, client *rpc.Client) *SanctumHandler {
	return &SanctumHandler{cfg: cfg, client: client}
}
