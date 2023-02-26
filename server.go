package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ronanzindev/go-cacher-example/cache"
)

type ServerOpts struct {
	ListenAddr  string
	IsLeader    bool
	LeaderAdder string
}

type Server struct {
	ServerOpts
	followers map[net.Conn]struct{} // followers são as conexões
	cache     cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		// TODO: only allocate this when we are the lead
		followers: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("server is running on port [%s]\n ", s.ListenAddr)

	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAdder)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("connected with the leader: ", s.LeaderAdder)
			s.handleConnection(conn)
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConnection(conn)

	}
}
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	if s.IsLeader {
		s.followers[conn] = struct{}{}
	}
	fmt.Println("connection made: ", conn.RemoteAddr())
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("connection read error: %s\n", err)
			break
		}

		go s.handleCommand(conn, buf[:n]) // buf[:n] é a mensagem, ex :  Set foo bar 1235
	}

}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)

	if err != nil {
		fmt.Println("failed to parse command:", err)
		conn.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("receiveid command %s\n", msg.Cmd)

	switch msg.Cmd {

	case CMDSet:
		err = s.handleSetCmd(conn, msg)

	case CMDGet:
		err = s.handleGetCmd(conn, msg)
	}

	if err != nil {
		fmt.Println("failed to handle command:", err)
		conn.Write([]byte(err.Error()))
	}

}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}
	go s.sendToFollowers(context.TODO(), msg)
	return nil
}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}

	_, err = conn.Write(val)

	return err

}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		fmt.Println("forwarding key to followers")
		rawMsg := msg.ToBytes()
		fmt.Println("forwarding key to followers: ", string(rawMsg))
		_, err := conn.Write(rawMsg)
		if err != nil {
			fmt.Println("write to follower error: ", err)
			continue
		}
	}
	return nil

}
