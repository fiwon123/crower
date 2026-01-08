package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateFileIfNotExists(filePath string) {
	file, err := os.OpenFile(filePath,
		os.O_CREATE|os.O_RDWR,
		0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func CreateFolderIfNotExists(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	return nil
}

func DeleteFile(path string) error {
	err := os.RemoveAll(path)
	return err
}

func FilePathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func FolderPathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

// Consider path if has extension
func IsValidFilePath(path string) bool {
	if path == "" {
		return false
	}

	if path == "." {
		return false
	}

	if !strings.Contains(path, "/") && !strings.Contains(path, "\\") {
		return false
	}

	base := filepath.Base(path)
	splitted := strings.Split(base, ".")

	return len(splitted) == 2
}

func IsValidFolderPath(path string) bool {
	if path == "" {
		return false
	}

	if path == "." {
		return false
	}

	if !strings.Contains(path, "/") && !strings.Contains(path, "\\") {
		return false
	}

	return true
}
