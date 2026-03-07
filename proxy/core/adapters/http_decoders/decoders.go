package http_decoders

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/core/domain"
)

var (
	ErrEmptyPayload = errors.New("err payload is empty")
)

type HTTPDecoder struct {
}

func NewHTTPDecoder() *HTTPDecoder {
	return &HTTPDecoder{}
}

func (deco *HTTPDecoder) FromConn(conn net.Conn) (*domain.HttpPackage, error) {
	body, err := deco.readCoon(conn)

	if err != nil {
		return nil, err
	}

	log.Println(body)
	return deco.stringToHttpPackage(body)
}

func (deco *HTTPDecoder) stringToHttpPackage(payload string) (*domain.HttpPackage, error) {
	if len(payload) == 0 {
		log.Println("payload is empty")
		return &domain.HttpPackage{}, ErrEmptyPayload
	}

	var (
		headers    = make(map[string]string)
		statusLine = ""
		body       = ""
		onBody     = false
	)

	for i, line := range strings.Split(payload, "\r\n") {
		log.Println("(", i, ")", "line:", line)

		if i == 0 {
			statusLine = line
			continue
		}

		if len(line) == 0 {
			onBody = true
			log.Println("Switching to body")
			continue
		}

		if onBody {
			body += line + "\r\n" // TODO: Maybe this is not necessary
		} else {
			header := strings.Split(line, ":")
			headers[header[0]] = header[1]
		}

	}

	return &domain.HttpPackage{
		StatusLine: statusLine,
		Headers:    headers,
		Body:       body,
	}, nil

}

// TODO: PERF: What about working with bytes instead of casting to string?
func (deco *HTTPDecoder) readCoon(conn net.Conn) (string, error) {
	var (
		buf        = make([]byte, 0)
		tmp        = make([]byte, 256)
		totalBytes = 0
	)

	for {
		lastByteRead, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			log.Println(err)
			return "", errors.Join(
				fmt.Errorf("err while reading the connection"),
				err,
			)
		}
		buf = append(buf, tmp[:lastByteRead]...)
		totalBytes += lastByteRead

		if lastByteRead < len(tmp) || errors.Is(err, io.EOF) {
			log.Println("Messaged read")
			break
		}
	}

	log.Println("total bytes", totalBytes)
	return string(buf[:totalBytes]), nil
}
