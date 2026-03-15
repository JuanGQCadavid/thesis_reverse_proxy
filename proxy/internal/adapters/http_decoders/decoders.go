package http_decoders

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
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

func (deco *HTTPDecoder) FromHttpResponse(resp *http.Response) (*domain.HttpPackage, error) {
	var (
		result = &domain.HttpPackage{
			StatusLine: domain.HttpStatusLineMultipart{
				HttpVersion: resp.Proto,
				StatusCode:  resp.Status,
			},
			Headers: make(map[string]string),
		}
	)

	// Getting Headers
	for k, v := range resp.Header {
		value := v[0]
		if len(v) > 1 {
			for _, vv := range v[1:] {
				value += ", " + vv
			}
		}
		result.Headers[k] = value
	}

	// Getting Body
	body, err := deco.readCoon(resp.Body)

	if err != nil {
		return result, errors.Join(fmt.Errorf("err whilr trying to read the body"), err)
	}

	result.Body = body
	result.BodyBytes = []byte(body)

	//log.Printf("%+v\n", result)

	return result, nil
}

func (deco *HTTPDecoder) stringToHttpPackage(payload string) (*domain.HttpPackage, error) {
	if len(payload) == 0 {
		log.Println("payload is empty")
		return &domain.HttpPackage{}, ErrEmptyPayload
	}

	var (
		headers    = make(map[string]string)
		statusLine = domain.HttpStatusLineMultipart{}
		body       = ""
		onBody     = false
	)

	for i, line := range strings.Split(payload, "\r\n") {
		//log.Println("(", i, ")", "line:", line)

		if i == 0 {
			tempStatusLine := strings.Split(line, " ")
			statusLine.Method = tempStatusLine[0]
			statusLine.Resource = tempStatusLine[1]
			statusLine.HttpVersion = tempStatusLine[2]
			continue
		}

		if len(line) == 0 {
			onBody = true
			//log.Println("Switching to body")
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
		BodyBytes:  []byte(body),
	}, nil

}

// TODO: PERF: What about working with bytes instead of casting to string?
func (deco *HTTPDecoder) readCoon(conn io.Reader) (string, error) {
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
