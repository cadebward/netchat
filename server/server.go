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
	ln, err := net.Listen(s.Network, s.Port)
	if err != nil {
		log.Panic(err)
	}
	log.Println(s.Network, "network listening on port", s.Port)
	// handle all incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		childCtx, _ := context.WithCancel(ctx)
		go s.handleConnection(childCtx, conn)
	}
}

func readUsername(conn net.Conn) (string, error) {
	conn.Write([]byte("What is your username? "))
	username, err := readInput(conn)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	username, err := readUsername(conn)
	if err != nil {
		log.Fatal(err)
	}

	// setup a new client and broadcast to all current clients
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
		Body:     "Welcome to netchat! Press Enter when you want to send.",
	})
	childCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// listen for any messages sent by each client
	go client.Listen(ctx)
	for {
		select {
		case <-childCtx.Done():
			return
		case message := <-client.Data:
			if message.Err != nil {
				// client has diconnected, remove connection and announce to clients
				client.Close()
				s.broadcast(Message{
					Time:     time.Now(),
					Err:      nil,
					Username: "system",
					Body:     client.Username + " has disconnected.",
				})
				return
			}
			// client has sent a message, broadcast to all clients
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
