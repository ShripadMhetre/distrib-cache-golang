package main

import (
	"fmt"
	"log"
	"net"

	"github.com/shripadmhetre/distrib-cache-golang/cache"
)

type ServerOptions struct {
	ListenAddr string
}

type Server struct {
	ServerOptions
	cache cache.Cache
}

func NewServer(ops ServerOptions, c cache.Cache) *Server {
	return &Server{
		ServerOptions: ops,
		cache:         c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("Server started listening at [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("connection made:", conn.RemoteAddr())

	fmt.Println("connection closed:", conn.RemoteAddr())
}
