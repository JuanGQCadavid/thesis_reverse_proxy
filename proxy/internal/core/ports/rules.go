package ports

import (
	"context"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type Rule interface {
	Execute(ctx context.Context, conn net.Conn, userAgentRequest *domain.HttpPackage) error
}
