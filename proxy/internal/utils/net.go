package utils

import (
	"log"
	"net"
)

func CloseDial(ln net.Listener) {
	if err := ln.Close(); err != nil {
		log.Fatal("err while closing the listener ", err.Error())
	}
}
