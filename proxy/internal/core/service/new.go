package service

import (
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
	srv.Config = config

	for _, rule := range config.Config.RulesConfig {
		if rule.Cache != nil {
			rule.RuleType = domain.CACHE_RULE
		} else if rule.Buffer != nil {
			rule.RuleType = domain.BUFFER_RULE
		} else {
			rule.RuleType = domain.PROXY_RULE
		}
	}

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
