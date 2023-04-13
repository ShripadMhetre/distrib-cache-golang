package main

import (
	"flag"

	"github.com/shripadmhetre/distrib-cache-golang/cache"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address of server")

	flag.Parse()

	options := ServerOptions{
		ListenAddr: *listenAddr,
	}

	server := NewServer(options, cache.New())
	server.Start()
}
