package util

import (
	"os"
	"path/filepath"
)

func CurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}

func CacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".malguem", "templates")

	return cacheDir, nil
}
