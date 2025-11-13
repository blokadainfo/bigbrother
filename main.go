package main

import (
	"log"

	"github.com/blokadainfo/bigbrother/src/config"
	"github.com/blokadainfo/bigbrother/src/record"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config from environment: %v", err)
	}

	for _, sc := range c.Streams {
		go record.StartRecord(c.Recordings, sc) // Start recording in a new goroutine
	}

	select {} // Block forever, so all goroutines can run
}
