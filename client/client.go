package client

import (
	"log"
	"net"

	// TODO how do i make a helper function like this w/o making a module
	"github.com/cadebward/netchat/input"
)

// type Message struct {
// 	Time int64
// 	Body string
// }

type Client struct {
	Conn     net.Conn
	Username string
	Room     string
	Data     chan string
	// Data     chan []*Message
}

func NewClient(conn net.Conn, username string) *Client {
	return &Client{
		Conn:     conn,
		Username: username,
		Room:     "general",
		Data:     make(chan string),
		// Data:     make(chan []*Message),
	}
}

func (c *Client) Close() {
	log.Println("connection closed", c.Username)
}

func (c *Client) Send() {
	log.Println("send", c.Username)
	for {
		msg, err := input.ReadInput(c.Conn)
		if err != nil {
			// TODO how do i catch the exited tenlet error?
			// if error is due to exiting telent
			// send connection closed event
			log.Panic(err)
		}
		log.Println(c.Username, "is sending message:", msg)
		c.Data <- msg
		// TODO impl rooms
	}
}

func (c *Client) Join(room string) {
	log.Println("join", c.Username)
}
