package remote

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"malguem/internal/util"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Clone(url string) (string, error) {
	template := extractRepoName(url)
	cacheDir, err := util.CacheDir()
	if err != nil {
		return "", err
	}
	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s-%s", template, hashRepoUrl(url)))

	// Make sure the cache directory exists
	os.MkdirAll(cacheDir, os.ModePerm)

	// Check if the template already cached
	if _, err := os.Stat(cachePath); err == nil {
		fmt.Printf("⬇️  Pulling updates for template: %s\n", cachePath)

		branch, err := getDefaultBranch(cachePath)
		if err != nil {
			return "", err
		}

		cmd := exec.Command("git", "pull", "origin", branch)
		cmd.Dir = cachePath

		err = cmd.Run()
		if err != nil {
			return "", err
		}

		return cachePath, nil
	}

	// Clone the template
	fmt.Printf("⏳  Cloning template into %s\n", cachePath)
	cmd := exec.Command("git", "clone", url, cachePath)
	cmd.Dir = cacheDir

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone template: %v\n", err)
	}

	return cachePath, nil
}

func CompareCommitHash(path, url string) bool {
	return false
}

func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, ".git")
}

func hashRepoUrl(url string) string {
	h := sha1.New()
	h.Write([]byte(url))

	return hex.EncodeToString(h.Sum(nil))[:8]
}

func getDefaultBranch(path string) (string, error) {
	cmd := exec.Command("git", "remote", "show", "origin")
	cmd.Dir = path

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get default branch: %w\nOutput: %s\n", err, output.String())
	}

	lines := strings.SplitSeq(output.String(), "\n")
	for line := range lines {
		line = strings.TrimSpace(line)

		if after, ok := strings.CutPrefix(line, "HEAD branch:"); ok {
			return strings.TrimSpace(after), nil
		}
	}

	return "", fmt.Errorf("default branch not found in output:\n%s\n", output.String())
}
