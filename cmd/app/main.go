package main

import (
	"log"
	"os"

	discordBot "github.com/evmartinelli/go-discordbot-panela/internal/app"
)

func main() {
	if err := discordBot.RunServer(); err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	return
}
