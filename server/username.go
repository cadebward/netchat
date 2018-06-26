package server

import (
	"net"

	"github.com/cadebward/netchat/input"
)

func readUsername(conn net.Conn) (string, error) {
	conn.Write([]byte("What is your username? "))
	username, err := input.ReadInput(conn)
	if err != nil {
		return "", err
	}
	return username, nil
}
