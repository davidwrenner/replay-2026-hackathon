package main

import (
	"log"
	"os"

	"github.com/davidwrenner/replay-2026-hackathon/api"
)

func main() {
	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := api.NewServer(addr)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
