package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/http_decoders"
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

var (
	defaultDecoder = http_decoders.HTTPDecoder{}
)

func CloseDial(ln net.Listener) {
	if err := ln.Close(); err != nil {
		log.Fatal("err while closing the listener ", err.Error())
	}
}

func ForwardConnection(req *domain.HttpPackage, originServer string) (*domain.HttpPackage, error) {
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
	defer osResponse.Body.Close()

	return defaultDecoder.FromHttpResponse(osResponse)
}
