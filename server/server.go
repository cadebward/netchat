package server

import (
	"io"
	"log"
	"net"
	"os"

	c "github.com/cadebward/netchat/client"
)

type Server struct {
	Port    string
	Network string
	Logfile string
	Clients []*c.Client
}

func NewServer(network string, port string, logfile string) *Server {
	return &Server{
		Network: network,
		Port:    port,
		Logfile: logfile,
		Clients: []*c.Client{},
	}
}

func (s *Server) Run() error {
	// setup logging
	outfile, err := os.Create(s.Logfile)
	if err != nil {
		log.Panic(err)
	}
	write := io.MultiWriter(os.Stdout, outfile)
	log.SetOutput(write)

	// start listening
	ln, err := net.Listen(s.Network, ":"+s.Port)
	if err != nil {
		return err
	}
	log.Println(s.Network, "network listening on port", s.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// func handleConnection(conn net.Conn) {
	username, err := readUsername(conn)
	if err != nil {
		log.Fatal(err)
	}

	client := c.NewClient(conn, username)
	// TODO lock a thing??
	s.Clients = append(s.Clients, client)
	client.Send()
	for {
		msg := <-client.Data
		log.Println("WHY U NO WORK!!??", msg)
	}
}
