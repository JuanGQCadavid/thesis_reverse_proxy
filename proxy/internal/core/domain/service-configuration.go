package domain

type RuleType uint8

const (
	CACHE_RULE RuleType = iota
	PROXY_RULE
	BUFFER_RULE
)

type ServiceConfiguration struct {
	Config GatewayConfig `json:"config"`
}

type GatewayConfig struct {
	Downstream  DownstreamConfig `json:"downstream"`
	RulesConfig []RulesConfig    `json:"rules"`
}

type RulesConfig struct {
	Regex          string    `json:"regex"`
	Name           string    `json:"name"`
	CollectMetrics bool      `json:"collect_metrics"`
	Cache          *struct{} `json:"cache"`  // HOLDER
	Buffer         *struct{} `json:"buffer"` // HOLDER

	// Internals
	RuleType RuleType
}

type DownstreamConfig struct {
	Principal PrincipalConfig `json:"principal"`
}

type PrincipalConfig struct {
	OriginServerURI string `json:"origin_server_uri"`
}
