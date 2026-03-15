package proxy

import (
	"context"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type ProxyRule struct {
	config *domain.ServiceConfiguration
}

func New(config *domain.ServiceConfiguration) *ProxyRule {
	return &ProxyRule{
		config: config,
	}
}

func (r ProxyRule) Execute(ctx context.Context, conn net.Conn, userAgentRequest *domain.HttpPackage) error {
	return nil
}
