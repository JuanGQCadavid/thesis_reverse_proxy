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

var (
	PathNotInConfig []byte
)

func init() {
	PathNotInConfig = domain.NewHttpPackage().
		WithMultipart(domain.HttpStatusLineMultipart{
			HttpVersion: "HTTP/1.1",
			StatusCode:  "404 Not Found",
		}).
		ToBytes()
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
			continue
		}
		go srv.handleConnV2(conn)
	}

	return nil
}

func buildResponse() string {
	return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
}

func (srv *Service) checkpointV2(conn net.Conn) {
	log.Println("dude?")
	if _, err := conn.Write([]byte(buildResponse())); err != nil {
		log.Println("err while writing data", err.Error())
	}

	log.Println("ya?")
	return
}

func (srv *Service) handleConnV2(conn net.Conn) {
	defer conn.Close() // Connection with user agent

	userAgentRequest, err := srv.decoder.FromConn(conn)

	if err != nil {
		log.Println("err while decoding connection", err.Error())
		return
	}

	var (
		ruleExecuted bool  = false
		ruleError    error = nil
	)

	log.Printf("%+v\n", srv.Config)

	for _, rule := range srv.Config.Config.RulesConfig {
		if match, _ := regexp.Match(rule.Regex, []byte(userAgentRequest.StatusLine.Resource)); match {
			log.Println("match", rule.Regex)
			log.Printf("%+v\n", rule)
			ruleError = srv.RulesMap[rule.RuleType].Execute(context.Background(), conn, userAgentRequest, &rule)
			ruleExecuted = true
			break
		}
	}
	//// TODO: What should I do
	//// DLQ?
	//// error to user?
	if ruleError != nil {
		log.Println("err while executing rule", ruleError.Error())
		conn.Write(
			domain.NewHttpPackage().
				WithMultipart(domain.HttpStatusLineMultipart{
					HttpVersion: "HTTP/1.1",
					StatusCode:  "500 Internal Server Error",
				}).
				WithBody(ruleError.Error()).
				ToBytes(),
		)
		return
	}

	if !ruleExecuted {
		log.Printf("%s request does not match any rule", userAgentRequest.StatusLine.Resource)
		conn.Write(PathNotInConfig)
		return
	}
}
