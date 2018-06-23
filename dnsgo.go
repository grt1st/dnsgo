package main

import (
	"github.com/grt1st/dnsgo/server"
	"time"
)

func main() {
	server := &server.Server{
		Host:     "127.0.0.1",
		Port:     53,
		RTimeout: 5 * time.Second,
		WTimeout: 5 * time.Second,
	}

	server.Run()
}
