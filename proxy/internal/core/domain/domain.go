package domain

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
)

type HttpPackage struct {
	StatusLine HttpStatusLineMultipart
	Headers    map[string]string
	Body       string
	BodyBytes  []byte
}

func NewHttpPackage() *HttpPackage {
	return &HttpPackage{
		Headers: make(map[string]string),
	}
}
func (p *HttpPackage) WithMultipart(StatusLine HttpStatusLineMultipart) *HttpPackage {
	p.StatusLine = StatusLine
	return p
}
func (p *HttpPackage) WithBody(body string) *HttpPackage {
	p.Body = body
	p.BodyBytes = []byte(body)
	return p
}

type HttpStatusLineMultipart struct {
	// Standard
	HttpVersion string

	// Request
	Resource string
	Method   string

	// Response
	StatusCode string
}

func (lm *HttpStatusLineMultipart) ToString() string {
	if len(lm.StatusCode) == 0 {
		return fmt.Sprintf("%s %s %s\r\n", lm.Method, lm.Resource, lm.HttpVersion)
	}
	return fmt.Sprintf("%s %s\r\n", lm.HttpVersion, lm.StatusCode)
}

type HttpBodyEncryption uint

const (
	Gzip HttpBodyEncryption = iota
)

type HttpConnectionType uint

const (
	KeepAlive HttpConnectionType = iota
)

func (p *HttpPackage) WithHost(host string) *HttpPackage {
	p.Headers["Host"] = host
	return p
}

func (p *HttpPackage) WithConnection(conType HttpConnectionType) *HttpPackage {
	switch conType {
	case KeepAlive:
		p.Headers["Connection"] = "keep-alive"
	}
	return p
}

func (p *HttpPackage) WithBodyEncryption(encryptionType HttpBodyEncryption) *HttpPackage {
	switch encryptionType {
	case Gzip:
		log.Println("Encrypting using gzip")
		p.Headers["Content-Encoding"] = "gzip"

		body_buffer := &bytes.Buffer{}
		writter := gzip.NewWriter(body_buffer)
		defer writter.Close()

		writter.Write([]byte(p.Body))
		writter.Flush()

		p.BodyBytes = body_buffer.Bytes()
		p.Headers["Content-Length"] = fmt.Sprintf("%d", len(p.BodyBytes))
	}

	return p
}

func (p *HttpPackage) headersToString() string {
	var headers string
	for k, v := range p.Headers {
		headers += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	return headers
}

func (p *HttpPackage) ToBytes() []byte {
	var (
		statusLine = p.StatusLine.ToString()
		headers    = p.headersToString()
		body       = p.Body
	)

	if len(p.BodyBytes) == 0 {
		return []byte(fmt.Sprintf("%s%s\r\n", statusLine, headers))
	}
	return []byte(fmt.Sprintf("%s%s\r\n%s ", statusLine, headers, body))
}
