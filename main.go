package main

import (
	"log"
	"time"

	"github.com/blokadainfo/bigbrother/src/autodelete"
	"github.com/blokadainfo/bigbrother/src/config"
	"github.com/blokadainfo/bigbrother/src/record"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config from environment: %v", err)
	}

	go autodelete.DeleteOldestFilesCron(c.Recordings.Path, c.Recordings.Limit, c.Recordings.DryRun, 5*time.Minute)

	for _, sc := range c.Streams {
		go record.StartRecord(c.Recordings, sc)
	}

	select {} // Block forever, so all goroutines can run
}
