package server

import (
	"bufio"
	"net"
	"strings"
)

func readInput(conn net.Conn) (string, error) {
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(input, "\r\n"), nil
}
