package util

import (
	"fmt"
	"io"
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

func CopyDir(source, destination string) error {
	src := filepath.Clean(source)
	dst := filepath.Clean(destination)

	// Check if source directory exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if source is not directory
	if !srcInfo.IsDir() {
		return fmt.Errorf("Source path is not directory!")
	}

	// Create identical directory in the destination
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relativePath)

		// Get entry info
		info, err := entry.Info()
		if err != nil {
			return err
		}

		switch {
		case entry.IsDir():
			return os.MkdirAll(dstPath, info.Mode())
		case info.Mode().IsRegular():
			return CopyFile(path, dstPath, info.Mode())
		default:
			return nil
		}
	})
}

func CopyFile(source, destination string, mode os.FileMode) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	outputFile, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, sourceFile); err != nil {
		return nil
	}

	return outputFile.Sync()
}
