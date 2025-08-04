package config

import (
	"os"
)

type Config struct {
	StorageType     string
	LocalStorageDir string
	TempDir        string
	FFmpegPath     string
}

func New() *Config {
	return &Config{
		StorageType:     getEnvOrDefault("STORAGE_TYPE", "local"),
		LocalStorageDir: getEnvOrDefault("LOCAL_STORAGE_DIR", "storage"),
		TempDir:        getEnvOrDefault("TEMP_DIR", "tmp"),
		FFmpegPath:     getEnvOrDefault("FFMPEG_PATH", "ffmpeg"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) EnsureDirectories() error {
	dirs := []string{c.LocalStorageDir, c.TempDir}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	
	return nil
}
