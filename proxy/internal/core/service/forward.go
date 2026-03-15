package service

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

func (srv *Service) ruleProxy(conn net.Conn, userAgentRequest *domain.HttpPackage) {
	originServerResponse, err := srv.forwardConnection(
		userAgentRequest,
		srv.Config.Config.Downstream.Principal.OriginServerURI,
	)

	if err != nil {
		log.Println("err while forwarding connection", err.Error())
	}

	if _, err := conn.Write(originServerResponse.ToBytes()); err != nil {
		log.Fatalln("err while writing data", err.Error())
	}
}

// TODO: How can we reuse OS connection
func (srv *Service) forwardConnection(req *domain.HttpPackage, originServer string) (*domain.HttpPackage, error) {
	//Configuring Proxy headers
	req.
		WithHost(srv.Config.Config.Downstream.Principal.OriginServerURI).
		WithConnection(domain.KeepAlive).
		WithBodyEncryption(domain.Gzip)

	var (
		osUri = fmt.Sprintf("http://%s%s", originServer, req.StatusLine.Resource)
	)

	osRequest, err := http.NewRequest(req.StatusLine.Method, osUri, bytes.NewBuffer(req.BodyBytes))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error creating http request"), err)
	}

	for k, v := range req.Headers {
		osRequest.Header.Add(k, v)
	}
	osResponse, err := http.DefaultClient.Do(osRequest)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error forwarding http request"), err)
	}

	return srv.decoder.FromHttpResponse(osResponse)
}
