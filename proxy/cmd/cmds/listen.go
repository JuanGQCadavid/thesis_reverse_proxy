package cmds

import (
	"log"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/cache/memory"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/http_decoders"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/rules/buffer"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/rules/cache"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/rules/proxy"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/service"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/utils"
	"github.com/spf13/cobra"
)

var (
	// Flags
	configFile string

	// Commands
	ListenCMD = &cobra.Command{
		Use:   "listen [ARGS]",
		Short: "Set up a proxy server",
		Long:  "Set up a proxy server with the configuration passed down by args",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(configFile)
		},
	}

	// Adapters
	memCache *memory.Cache = memory.NewCache()
)

func init() {
	ListenCMD.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "config yaml file")
}

func execute(configFilePath string) error {
	var (
		config *domain.ServiceConfiguration = &domain.ServiceConfiguration{}

		// Rules
		proxyRule = proxy.New(config)
		cacheRule = cache.New().
				WithRepository(memCache).
				WithConfig(config).
				WithProxyRule(*proxyRule)
	)
	if err := utils.FromFilePathToStruct(configFilePath, config); err != nil {
		log.Fatal("err to read configuration file ", err.Error())
	}

	log.Println("load configuration successfully")
	log.Printf("%+v", config)

	return service.New().
		WithConfig(config).
		WithDecoder(http_decoders.HTTPDecoder{}).
		WithRule(domain.PROXY_RULE, proxyRule).
		WithRule(domain.BUFFER_RULE, buffer.New()).
		WithRule(domain.CACHE_RULE, cacheRule).
		Run()
}
