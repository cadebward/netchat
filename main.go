package main

import (
	"encoding/json"
	"log"
	"os"

	s "github.com/cadebward/netchat/server"
)

type Configuration struct {
	Network string
	Port    string
	Logfile string
}

func readConfig() (Configuration, error) {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	decoder.Decode(&config)
	return config, nil
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Panic(err)
	}
	server := s.NewServer(config.Network, config.Port, config.Logfile)
	server.Run()
}
