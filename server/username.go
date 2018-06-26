package server

import (
	"net"
)

func readUsername(conn net.Conn) (string, error) {
	conn.Write([]byte("What is your username? "))
	username, err := readInput(conn)
	if err != nil {
		return "", err
	}
	return username, nil
}
