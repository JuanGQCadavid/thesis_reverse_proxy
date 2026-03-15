package cache

import (
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/rules/proxy"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/ports"
)

func New() *Rule {
	return &Rule{}
}

func (r *Rule) WithConfig(config *domain.ServiceConfiguration) *Rule {
	r.config = config
	return r
}

func (r *Rule) WithProxyRule(proxyRule proxy.Rule) *Rule {
	r.proxyRule = proxyRule
	return r
}

func (r *Rule) WithRepository(repo ports.Cache) *Rule {
	r.repo = repo
	return r
}
