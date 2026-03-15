package domain

type RuleType string

const (
	CACHE_RULE  RuleType = "CACHE_RULE"
	PROXY_RULE  RuleType = "PROXY_RULE"
	BUFFER_RULE RuleType = "BUFFER_RULE"
)

type ServiceConfiguration struct {
	Config GatewayConfig `json:"config"`
}

type GatewayConfig struct {
	Downstream  DownstreamConfig `json:"downstream"`
	RulesConfig []RulesConfig    `json:"rules"`
}

type RulesConfig struct {
	Regex          string       `json:"regex"`
	Name           string       `json:"name"`
	CollectMetrics bool         `json:"collect_metrics"`
	Cache          CacheConfig  `json:"cache"`  // HOLDER
	Buffer         BufferConfig `json:"buffer"` // HOLDER
	Proxy          ProxyConfig  `json:"proxy"`

	// Internals
	RuleType RuleType
}

type CacheConfig struct {
	Enable bool   `json:"enable"`
	TTL    string `json:"ttl"`
}

type BufferConfig struct {
	Enable bool `json:"enable"`
}

type ProxyConfig struct {
	Compress CompressConfig `json:"compress"`
}

type CompressConfig struct {
	Enable         bool   `json:"enable"`
	CompressMethod string `json:"method"`
}

type DownstreamConfig struct {
	Principal PrincipalConfig `json:"principal"`
}

type PrincipalConfig struct {
	OriginServerURI string `json:"origin_server_uri"`
}
