package remote

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Clone(template, url string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".malguem", "templates")
	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s-%s", template, hashRepoUrl(url)))

	// Make sure the cache directory exists
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check iff the template already cached
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, fmt.Errorf("template `%s` already fetched\n", url)
	}

	// Clone the template
	cmd := exec.Command("git", "clone", url, cachePath)
	cmd.Dir = cacheDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone template: %v\n", err)
	}

	return cachePath, nil
}

func hashRepoUrl(url string) string {
	h := sha1.New()
	h.Write([]byte(url))

	return hex.EncodeToString(h.Sum(nil))[:8]
}
