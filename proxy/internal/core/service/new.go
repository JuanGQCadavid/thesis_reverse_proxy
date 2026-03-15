package service

import (
	"log"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/http_decoders"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/ports"
)

func New() *Service {
	return &Service{
		RulesMap: make(map[domain.RuleType]ports.Rule),
	}
}

func (srv *Service) WithConfig(config *domain.ServiceConfiguration) *Service {
	log.Println("Configuring rules")

	for i, rule := range config.Config.RulesConfig {
		if rule.Cache.Enable {
			log.Println("rule ", rule, " is cache")
			rule.RuleType = domain.CACHE_RULE
		} else if rule.Buffer.Enable {
			log.Println("rule ", rule, " is buffer")
			rule.RuleType = domain.BUFFER_RULE
		} else {
			log.Println("rule ", rule, " is proxy")
			rule.RuleType = domain.PROXY_RULE
		}
		config.Config.RulesConfig[i] = rule
	}
	srv.Config = config

	return srv
}

func (srv *Service) WithDecoder(decoder http_decoders.HTTPDecoder) *Service {
	srv.decoder = decoder
	return srv
}

func (srv *Service) WithRule(ruleType domain.RuleType, rule ports.Rule) *Service {
	srv.RulesMap[ruleType] = rule
	return srv
}
