package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/shripadmhetre/distrib-cache-golang/cache"
	"github.com/shripadmhetre/distrib-cache-golang/client"
	"github.com/shripadmhetre/distrib-cache-golang/service"
)

type ServerOptions struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOptions
	cache   cache.Cache
	clients map[*client.Client]struct{}
}

func NewServer(ops ServerOptions, c cache.Cache) *Server {
	return &Server{
		ServerOptions: ops,
		cache:         c,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return fmt.Errorf("server listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println("Server.Run() => Error while dialing leader: ", err)
			}
		}()
	}

	log.Printf("Server started listening at [%s]\n", s.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}

	log.Println("connected to leader:", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, service.CmdJoin)

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	//fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := service.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}

	// fmt.Println("connection closed:", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch cmdType := cmd.(type) {
	case *service.CommandSet:
		s.handleSetCommand(conn, cmdType)
	case *service.CommandGet:
		s.handleGetCommand(conn, cmdType)
	case *service.CommandExists:
		s.handleExistsCommand(conn, cmdType)
	case *service.CommandJoin:
		s.handleJoinCommand(conn, cmdType)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *service.CommandJoin) error {
	fmt.Println("member joined the cluster:", conn.RemoteAddr())

	s.clients[client.NewClient(conn)] = struct{}{}

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *service.CommandGet) error {
	log.Printf("GET %s", cmd.Key)

	resp := service.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = service.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = service.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *service.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)

	go func() {
		for client := range s.clients {
			err := client.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	resp := service.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = service.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = service.StatusOK
	_, err := conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleExistsCommand(conn net.Conn, cmd *service.CommandExists) error {
	log.Printf("Exists %s", cmd.Key)

	resp := service.ResponseExists{}
	ok, err := s.cache.Exists(cmd.Key)

	if err != nil {
		resp.Status = service.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	if !ok {
		resp.Status = service.StatusKeyNotFound
	} else {
		resp.Status = service.StatusOK
	}
	_, err = conn.Write(resp.Bytes())

	return err
}
