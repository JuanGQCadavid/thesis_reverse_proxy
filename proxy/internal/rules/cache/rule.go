package cache

import (
	"context"
	"log"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type Rule struct {
}

func New() *Rule {
	return &Rule{}
}

func (r *Rule) Execute(
	ctx context.Context,
	conn net.Conn,
	userAgentRequest *domain.HttpPackage,
	ruleConfig *domain.RulesConfig,
) error {
	log.Println("Rule Cache", ruleConfig)
	return nil
}
