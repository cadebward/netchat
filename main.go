package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// TODO handle client disconnect
// TODO finish logging all events

type Configuration struct {
	Port string
}

type Client struct {
	Conn     net.Conn
	Username string
	Room     string
	Message  chan string
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var (
	outfile, _ = os.Create("./server.log")
	write      = io.MultiWriter(os.Stdout, outfile)
	clients    = map[string]*Client{}
)

func main() {
	log.SetOutput(write)
	// decode config file
	config, _ := readConfig()
	// start listening on specified port
	listener, _ := net.Listen("tcp", ":"+config.Port)
	log.Println("server started", config.Port)
	for {
		conn, err := listener.Accept()
		check(err)
		go handleConnection(conn)
	}
}

func readConfig() (Configuration, error) {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	decoder.Decode(&config)
	return config, nil
}

func getUsername(conn net.Conn) (string, error) {
	username, err := readInput(conn, "Welcome! What is your username? ")
	check(err)
	if clients[username] != nil {
		writeMessage(conn, "System", "That username is already taken!")
		return getUsername(conn)
	}
	return username, err
}

func handleConnection(conn net.Conn) {
	username, err := getUsername(conn)
	check(err)
	client := &Client{
		Conn:     conn,
		Username: username,
		Room:     "general",
		Message:  make(chan string),
	}
	log.Printf("connection established from %v as %v", client.Conn.RemoteAddr(), client.Username)
	// TODO remove clients when they disconnect
	clients[username] = client
	writeMessage(conn, "System", "Welcome, "+client.Username+"!\n")
	Broadcast(client.Username + " has joined the server!")
	// fire off routines for newly connected client
	go client.Send()
	// go client.Send()
	// go client.Receive()
}

func Broadcast(msg string) {
	for _, client := range clients {
		writeMessage(client.Conn, "System", msg+"\n")
	}
}

func (c *Client) Send() {
	for {
		msg, err := readInput(c.Conn, "")
		check(err)
		writeOwnMessage(c.Conn, c.Username, msg)
		for k, v := range clients {
			if k != c.Username {
				writeMessage(v.Conn, c.Username, msg)
			}
		}
	}
}

func readInput(conn net.Conn, output string) (string, error) {
	if output != "" {
		writeMessage(conn, "System", output)
	}
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("error from user readinput %v", err)
		return "", err
	}
	input = strings.Trim(input, "\r\n")
	return input, nil
}

// figure out a better way to handle broadcasting message
// TODO remove duplication between these two functions
func writeMessage(conn net.Conn, username, output string) {
	ts := time.Now().Format("02/01/2006 15:04:05")
	output = strings.Trim(output, "\r\n")
	msg := ts + " (" + username + "): " + output + "\n"
	conn.Write([]byte(msg))
}

func writeOwnMessage(conn net.Conn, username, output string) {
	ts := time.Now().Format("02/01/2006 15:04:05")
	output = strings.Trim(output, "\r\n")
	msg := "\033[1A"
	msg += "\033[K"
	msg += ts + " (" + username + "): " + output + "\n"
	conn.Write([]byte(msg))
}
