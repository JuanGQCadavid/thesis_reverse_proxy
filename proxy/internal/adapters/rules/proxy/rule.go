package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/utils"
)

type Rule struct {
	config *domain.ServiceConfiguration
}

func New(config *domain.ServiceConfiguration) *Rule {
	return &Rule{
		config: config,
	}
}

func (r *Rule) Execute(
	ctx context.Context,
	conn net.Conn,
	userAgentRequest *domain.HttpPackage,
	ruleConfig *domain.RulesConfig,
) error {
	log.Println("Executing rule", ruleConfig.RuleType)

	originServerResponse, err := utils.ProxyRequest(
		ctx,
		r.config.Config.Downstream.Principal.OriginServerURI,
		userAgentRequest,
		ruleConfig,
	)

	if err != nil {
		return errors.Join(err, fmt.Errorf("err while forwarding connection"))
	}

	log.Printf("%+v\n", originServerResponse)

	if _, err := conn.Write(originServerResponse.ToBytes()); err != nil {
		return errors.Join(err, fmt.Errorf("err while writing data"))
	}
	log.Println(" No error occurred")

	return nil
}
