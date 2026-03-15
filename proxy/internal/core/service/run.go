package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/http_decoders"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/ports"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/utils"
)

type Service struct {
	Config   *domain.ServiceConfiguration
	RulesMap map[domain.RuleType]ports.Rule
	decoder  http_decoders.HTTPDecoder
}

func (srv *Service) Run() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return errors.Join(fmt.Errorf("error listening on port 8080"), err)
	}
	defer utils.CloseDial(ln)

	// TODO: I should listen here control C to abort the process
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("err while accepting connection", err.Error())
		}
		go srv.handleConnV2(conn)
	}

	return nil
}

func (srv *Service) handleConnV2(conn net.Conn) {
	defer conn.Close() // Connection with user agent

	userAgentRequest, err := srv.decoder.FromConn(conn)
	if err != nil {
		log.Println("err while decoding connection", err.Error())
		return
	}

	for _, rule := range srv.Config.Config.RulesConfig {
		// TODO: Evaluate regex at init
		if match, _ := regexp.Match(rule.Regex, []byte(userAgentRequest.StatusLine.Resource)); match {
			if err = srv.RulesMap[rule.RuleType].Execute(context.Background(), conn, userAgentRequest); err != nil {
				log.Println("err while executing rule", err.Error())
			}
		}
	}
}
