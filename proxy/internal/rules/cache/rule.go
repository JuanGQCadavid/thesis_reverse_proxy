package cache

import (
	"context"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type CacheRule struct {
}

func New() *CacheRule {
	return &CacheRule{}
}

func (r CacheRule) Execute(ctx context.Context, conn net.Conn, userAgentRequest *domain.HttpPackage) error {
	return nil
}
