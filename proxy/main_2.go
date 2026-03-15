package main

import (
	"log"
	"net"

	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/adapters/http_decoders"
)

func init() {}

var (
	forwardTo = "localhost:9000"
	deco      http_decoders.HTTPDecoder
)

func init() {
	deco = http_decoders.HTTPDecoder{}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("err while setting up Listener ", err.Error())
	}
	defer closeDial(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("err while accepting connection", err.Error())
		}

		handleConnV2(conn)
	}

}

func handleConnV2(conn net.Conn) {
	defer conn.Close()

	resp, err := deco.FromConn(conn)
	if err != nil {
		return
	}

	log.Println(resp.StatusLine)
	log.Println(len(resp.Headers), resp.Headers)
	log.Println(resp.Body)

	if _, err := conn.Write([]byte(buildResponse())); err != nil {
		log.Fatalln("err while writing data", err.Error())
	}
}

func buildResponse() string {
	return "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
}

func closeDial(ln net.Listener) {
	if err := ln.Close(); err != nil {
		log.Fatal("err while closing the listener ", err.Error())
	}
}
