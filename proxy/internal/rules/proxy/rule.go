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

	userAgentRequest.
		WithHost(r.config.Config.Downstream.Principal.OriginServerURI).
		WithConnection(domain.KeepAlive)

	if ruleConfig.Proxy.Compress.Enable {
		userAgentRequest.WithBodyEncryption(domain.Gzip)
	}

	originServerResponse, err := utils.ForwardConnection(
		userAgentRequest,
		r.config.Config.Downstream.Principal.OriginServerURI,
	)

	if err != nil {
		return errors.Join(err, fmt.Errorf("err while forwarding connection"))
	}

	log.Println(" No error occurred")

	log.Printf("%+v\n", originServerResponse)

	if _, err := conn.Write(originServerResponse.ToBytes()); err != nil {
		return errors.Join(err, fmt.Errorf("err while writing data"))
	}
	log.Println(" No error occurred")

	return nil
}
