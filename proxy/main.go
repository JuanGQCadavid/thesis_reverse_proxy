package main

import (
	"io"
	"log"
	"net"
)

func init() {}

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

		handelConn(conn)
	}

}

func handelConn(conn net.Conn) {
	defer conn.Close()

	log.Println("handelConn start")
	defer log.Println("handelConn end")

	log.Println(conn.RemoteAddr())
	log.Println(conn.LocalAddr())

	buf := make([]byte, 0)
	tmp := make([]byte, 256)

	for {
		n, err := conn.Read(tmp)

		if err != nil {
			if err != io.EOF {
				log.Fatalln("err while reading data", err.Error())
			}
			break
		}

		buf = append(buf, tmp[:n]...)

		if n < len(tmp) {
			break
		}
	}
	log.Println(len(buf))
	log.Println(string(buf))
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
