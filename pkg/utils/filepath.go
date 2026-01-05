package utils

import (
	"fmt"
	"log"
	"os"
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

func IsValidFilePath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func IsValidFolderPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}
