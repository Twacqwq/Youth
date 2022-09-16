package utils

import (
	"os"
	"path/filepath"
)

func LoadConfig(configDir string) (*os.File, error) {
	configPath := filepath.Join(configDir, "config.json")
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}