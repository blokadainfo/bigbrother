package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/blokadainfo/bigbrother/src/config"
)

const (
	filePrefix  = "%Y-%m-%d_%H:%M:%S"
	segmentTime = 300 // 5 minutes
)

func CreateFFmpegCommand(rc config.RecordingsConfig, sc config.StreamConfig) *exec.Cmd {
	cmdArgs := []string{
		"-i", sc.Src,
		"-c", "copy",
		"-f", "segment",
		"-segment_time", strconv.Itoa(segmentTime),
		"-segment_format", "mkv",
		"-strftime", "1",
		"-reset_timestamps", "1",
		fmt.Sprintf("%s/%s_%s.mkv", rc.Path, filePrefix, sc.Dst),
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
