package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/thantko20/tubbym-transcoder/pkg/config"
	"github.com/thantko20/tubbym-transcoder/pkg/storage"
	"github.com/thantko20/tubbym-transcoder/pkg/transcoder"
)

func main() {
	// Initialize configuration
	cfg := config.New()
	if err := cfg.EnsureDirectories(); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	// Initialize storage
	var store storage.Storage
	switch cfg.StorageType {
	case "local":
		store = storage.NewLocalStorage(cfg.LocalStorageDir)
	default:
		log.Fatalf("Unsupported storage type: %s", cfg.StorageType)
	}

	// Initialize transcoder
	tr := transcoder.New(cfg.FFmpegPath, cfg.TempDir)

	// Example: Process a video file
	inputFile := os.Args[1] // Get input file from command line argument
	outputPath := filepath.Join("videos", filepath.Base(inputFile))

	// Download the input file to a temporary location
	tmpInput := filepath.Join(cfg.TempDir, "input", filepath.Base(inputFile))
	if err := os.MkdirAll(filepath.Dir(tmpInput), 0755); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	if err := store.Download(inputFile, tmpInput); err != nil {
		log.Fatalf("Failed to download input file: %v", err)
	}

	// Transcode the video
	outputDir, err := tr.TranscodeToHLS(tmpInput)
	if err != nil {
		log.Fatalf("Failed to transcode video: %v", err)
	}
	defer tr.Cleanup(outputDir)

	// Upload the transcoded files
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Calculate relative path for the destination
		relPath, err := filepath.Rel(outputDir, path)
		if err != nil {
			return err
		}

		// Upload the file
		dst := filepath.Join(outputPath, relPath)
		if err := store.Upload(path, dst); err != nil {
			return fmt.Errorf("failed to upload %s: %v", relPath, err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to upload transcoded files: %v", err)
	}

	fmt.Printf("Successfully transcoded %s to HLS format\n", inputFile)
}