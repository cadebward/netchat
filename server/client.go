package server

import (
	"context"
	"log"
	"net"
	"strings"
	"time"
)

type Message struct {
	Time     time.Time
	Body     string
	Err      error
	Username string
}

type Client struct {
	Conn     net.Conn
	Username string
	Room     string
	Data     chan Message
	isClosed bool
}

func NewClient(conn net.Conn, username string) *Client {
	return &Client{
		Conn:     conn,
		Username: username,
		Room:     "general",
		Data:     make(chan Message),
		isClosed: false,
	}
}

func (c *Client) Close() {
	// close out the connection and set a flag so our goroutine will stop
	log.Println(c.Username, " has disconnected")
	c.isClosed = true
	c.Conn.Close()
}

func (c *Client) Listen(ctx context.Context) {
	for {
		if c.isClosed {
			return
		}
		msg, err := readInput(c.Conn)
		log.Println(c.Username, "is sending message:", msg)
		c.Data <- Message{
			Time:     time.Now(),
			Body:     msg,
			Err:      err,
			Username: c.Username,
		}
	}
}

func (c *Client) Send(message Message) (int, error) {
	formatted := formatMessage(message)
	if c.Username == message.Username {
		// removes the chat message and replaces it with the newly formatted version
		// of the chat message when the client is the one who sent the message
		resetCursor := "\033[1A\033[K"
		return c.Conn.Write([]byte(resetCursor + formatted))
	} else {
		return c.Conn.Write([]byte(formatted))
	}
}

func formatMessage(message Message) string {
	ts := message.Time.Format("02/01/2006 15:04:05")
	output := strings.Trim(message.Body, "\r\n")
	msg := ts + " (" + message.Username + "): " + output + "\n"
	return msg
}
