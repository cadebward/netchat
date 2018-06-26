package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/cadebward/netchat/server"
)

func readConfig() (server.Configuration, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return server.Configuration{}, err
	}
	defer file.Close()
	config := server.Configuration{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return server.Configuration{}, err
	}
	return config, nil
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Panic(err)
	}
	s := server.NewServer(config)
	ctx := context.Background()
	s.Run(ctx)
}
