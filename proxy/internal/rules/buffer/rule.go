package buffer

import (
	"context"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type BufferRule struct {
}

func New() *BufferRule {
	return &BufferRule{}
}

func (r BufferRule) Execute(ctx context.Context, conn net.Conn, userAgentRequest *domain.HttpPackage) error {
	return nil
}
