package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Recordings RecordingsConfig
	Streams    []StreamConfig
}

type RecordingsConfig struct {
	Path  string
	Hours int
}

type StreamConfig struct {
	Src string
	Dst string
}

func LoadConfig() (Config, error) {
	recordingsConfig := RecordingsConfig{
		Path: os.Getenv("RECORDINGS_PATH"),
	}

	if recordingsConfig.Path == "" {
		return Config{}, fmt.Errorf("RECORDINGS_PATH is required")
	}

	recordingsHoursS := os.Getenv("RECORDINGS_HOURS")
	if recordingsHoursS == "" {
		return Config{}, fmt.Errorf("RECORDINGS_HOURS is required")
	}

	recordingsHours, err := strconv.Atoi(recordingsHoursS)
	if err != nil {
		return Config{}, fmt.Errorf("invalid RECORDINGS_HOURS value: %v", err)
	}
	recordingsConfig.Hours = recordingsHours

	var streams []StreamConfig
	for i := 0; ; i++ {
		src := os.Getenv(fmt.Sprintf("STREAMS_%d_SRC", i))
		dst := os.Getenv(fmt.Sprintf("STREAMS_%d_DST", i))
		if src == "" || dst == "" {
			break
		}
		streams = append(streams, StreamConfig{Src: src, Dst: dst})
	}

	if len(streams) == 0 {
		return Config{}, fmt.Errorf("Streams must be defined with STREAMS_[num]_SRC and STREAMS_[num]_DST (starting from STREAMS_0_SRC/DST)")
	}

	return Config{
		Recordings: recordingsConfig,
		Streams:    streams,
	}, nil
}
