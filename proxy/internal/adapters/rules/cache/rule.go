package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/rules/proxy"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/ports"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/utils"
)

var (
	AllowedMethods = []byte("/GET|OPTIONS/g")
	Seconds        = []byte("/[0-9]*s/gm")
	Minutes        = []byte("/[0-9]*m/gm")
	defaultTTL     = time.Duration(5 * time.Minute)
	cacheTTLSTypes = map[string]*time.Duration{}
)

type Rule struct {
	repo      ports.Cache
	proxyRule proxy.Rule
	config    *domain.ServiceConfiguration
}

func (r *Rule) Execute(
	ctx context.Context,
	conn net.Conn,
	userAgentRequest *domain.HttpPackage,
	ruleConfig *domain.RulesConfig,
) error {
	log.Println("Rule Cache", ruleConfig)

	if match, _ := regexp.Match(userAgentRequest.StatusLine.Method, AllowedMethods); !match {
		log.Printf("%s request does not match any allowd methods", userAgentRequest.StatusLine.Resource)
		return r.proxyRule.Execute(ctx, conn, userAgentRequest, ruleConfig)
	}

	data, err := r.repo.Get(userAgentRequest.StatusLine)

	if err != nil {
		log.Printf("Error getting data from repo: %v", err)
	}

	if data != nil {
		log.Println("Cache Hit")
		return r.writeToConn(conn, data.Data)
	}

	log.Println("Cache Miss")

	originServerResponse, err := utils.ProxyRequest(
		ctx,
		r.config.Config.Downstream.Principal.OriginServerURI,
		userAgentRequest,
		ruleConfig,
	)

	if err != nil {
		return errors.Join(err, fmt.Errorf("err while forwarding connection"))
	}

	if err = r.repo.Save(userAgentRequest.StatusLine, &domain.CacheData{
		Data: originServerResponse,
		TTL:  r.getTTL(ruleConfig),
	}); err != nil {
		log.Printf("Error caching response: %v\n", err)
	}

	return r.writeToConn(conn, originServerResponse)
}

func (r *Rule) getTTL(ruleConfig *domain.RulesConfig) time.Time {
	if cacheTTLSTypes[ruleConfig.Cache.TTL] != nil {
		log.Println("Returning from ttl cache")
		return time.Now().Add(*cacheTTLSTypes[ruleConfig.Cache.TTL])
	}

	var ttl time.Duration = defaultTTL

	if minutes, _ := regexp.Match(ruleConfig.Cache.TTL, Minutes); minutes {
		log.Println("Reading ttl on minutes")
		vals, err := strconv.Atoi(ruleConfig.Cache.TTL[0 : len(ruleConfig.Cache.TTL)-1])

		if err != nil {
			log.Printf("Error converting ttl %s to int: %v\n", ruleConfig.Cache.TTL, err)
		}
		ttl = time.Duration(vals) * time.Minute

	} else if seconds, _ := regexp.Match(ruleConfig.Cache.TTL, Seconds); seconds {
		log.Println("Reading ttl on seconds")
		vals, err := strconv.Atoi(ruleConfig.Cache.TTL[0 : len(ruleConfig.Cache.TTL)-1])

		if err != nil {
			log.Printf("Error converting ttl %s to int: %v\n", ruleConfig.Cache.TTL, err)
		}

		ttl = time.Duration(vals) * time.Second
	}
	cacheTTLSTypes[ruleConfig.Cache.TTL] = &ttl
	return time.Now().Add(ttl)
}

func (r *Rule) writeToConn(conn net.Conn, data *domain.HttpPackage) error {
	if _, err := conn.Write(data.ToBytes()); err != nil {
		return errors.Join(err, fmt.Errorf("err while writing data"))
	}
	return nil
}
