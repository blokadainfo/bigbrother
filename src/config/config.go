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
	Path   string
	Limit  int  // In gigabytes (GB), 0 means no limit
	DryRun bool // When set to true only prints what would get deleted instead of actually removing files
}

type StreamConfig struct {
	Src string
	Dst string
}

func LoadConfig() (Config, error) {
	recordingsPath, recordingsPathF := os.LookupEnv("RECORDINGS_PATH")
	if !recordingsPathF {
		return Config{}, fmt.Errorf("Env var RECORDINGS_PATH must be set")
	}
	if recordingsPath == "" {
		return Config{}, fmt.Errorf("Env var RECORDINGS_PATH must not be an empty string")
	}

	recordingsLimitS, recordingsLimitF := os.LookupEnv("RECORDINGS_LIMIT")
	if !recordingsLimitF {
		return Config{}, fmt.Errorf("Env var RECORDINGS_LIMIT must be set")
	}
	if recordingsLimitS == "" {
		return Config{}, fmt.Errorf("Env var RECORDINGS_LIMIT must not be an empty string")
	}
	recordingsLimit, err := strconv.Atoi(recordingsLimitS)
	if err != nil {
		return Config{}, fmt.Errorf("Invalid RECORDINGS_LIMIT value: %v", err)
	}

	recordingsDryRunS, recordingsDryRunF := os.LookupEnv("RECORDINGS_DRYRUN")
	if !recordingsDryRunF {
		return Config{}, fmt.Errorf("Env var RECORDINGS_DRYRUN must be set")
	}
	if recordingsDryRunS == "" {
		return Config{}, fmt.Errorf("Env var RECORDINGS_DRYRUN must not be an empty string")
	}
	recordingsDryRun, err := strconv.ParseBool(recordingsDryRunS)
	if err != nil {
		return Config{}, fmt.Errorf("Invalid RECORDINGS_DRYRUN value: %v", err)
	}

	var streams []StreamConfig
	for i := 0; ; i++ {
		src, srcF := os.LookupEnv(fmt.Sprintf("STREAMS_%d_SRC", i))
		dst, dstF := os.LookupEnv(fmt.Sprintf("STREAMS_%d_DST", i))

		// If neither source nor destination are set, all streams have been found
		if !srcF && !dstF {
			break
		}
		// If only one of the two isn't found, it's user error
		if !srcF || !dstF {
			return Config{}, fmt.Errorf("Both STREAMS_%d_SRC and STREAMS_%d_DST must be set", i, i)
		}

		if src == "" {
			return Config{}, fmt.Errorf("Env var STREAMS_%d_SRC must not be an empty string", i)
		}
		if dst == "" {
			return Config{}, fmt.Errorf("Env var STREAMS_%d_DST must not be an empty string", i)
		}

		streams = append(streams, StreamConfig{Src: src, Dst: dst})
	}

	if len(streams) == 0 {
		return Config{}, fmt.Errorf("Streams must be defined with STREAMS_[num]_SRC and STREAMS_[num]_DST (starting from STREAMS_0_SRC/DST)")
	}

	return Config{
		Recordings: RecordingsConfig{
			Path:   recordingsPath,
			Limit:  recordingsLimit,
			DryRun: recordingsDryRun,
		},
		Streams: streams,
	}, nil
}
