package transcoder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Transcoder struct {
	ffmpegPath string
	tempDir    string
}

func New(ffmpegPath, tempDir string) *Transcoder {
	return &Transcoder{
		ffmpegPath: ffmpegPath,
		tempDir:    tempDir,
	}
}

func (t *Transcoder) TranscodeToHLS(inputPath string) (string, error) {
	// Create a temporary output directory for the HLS stream
	outputDir := filepath.Join(t.tempDir, filepath.Base(inputPath)+"-hls")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	playlistPath := filepath.Join(outputDir, "playlist.m3u8")

	// FFmpeg command to transcode video to HLS
	cmd := exec.Command(t.ffmpegPath,
		"-i", inputPath,
		"-profile:v", "baseline", // baseline profile for better compatibility
		"-level", "3.0",
		"-start_number", "0",
		"-hls_time", "10", // 10 second segments
		"-hls_list_size", "0", // Keep all segments
		"-f", "hls",
		playlistPath,
	)

	// Capture stderr for logging
	cmd.Stderr = os.Stderr

	// Run the FFmpeg command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("transcoding failed: %w", err)
	}

	return outputDir, nil
}

func (t *Transcoder) Cleanup(outputDir string) error {
	return os.RemoveAll(outputDir)
}
