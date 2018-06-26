package server

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Configuration struct {
	Network string
	Port    string
	Logfile string
}

type Server struct {
	Port    string
	Network string
	Logfile string
	Clients []*Client
	mu      sync.Mutex
}

func NewServer(config Configuration) *Server {
	return &Server{
		Network: config.Network,
		Port:    config.Port,
		Logfile: config.Logfile,
		Clients: []*Client{},
		mu:      sync.Mutex{},
	}
}

func (s *Server) Run(ctx context.Context) error {
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
		log.Panic(err)
	}
	log.Println(s.Network, "network listening on port", s.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		childCtx, _ := context.WithCancel(ctx)
		go s.handleConnection(childCtx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	// func handleConnection(conn net.Conn) {
	username, err := readUsername(conn)
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(conn, username)
	s.mu.Lock()
	s.Clients = append(s.Clients, client)
	s.mu.Unlock()
	s.emit(Message{
		Time:     time.Now(),
		Err:      nil,
		Username: "system",
		Body:     client.Username + " has joined the server!",
	}, client.Username)
	client.Send(Message{
		Time:     time.Now(),
		Err:      nil,
		Username: "system",
		Body:     "Welcome to netchat! Feel free to type away. Press Enter when you want to send.",
	})
	childCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	go client.Listen(ctx)
	for {
		select {
		case <-childCtx.Done():
			return
		case message := <-client.Data:
			if message.Err != nil {
				client.Close()
				return
			}
			s.broadcast(message)
		}
	}
}

// Send a message to all connected clients
func (s *Server) broadcast(message Message) {
	for _, v := range s.Clients {
		v.Send(message)
	}
}

// Send a message to all clients minus one
func (s *Server) emit(message Message, username string) {
	for _, v := range s.Clients {
		if username != v.Username {
			v.Send(message)
		}
	}
}
