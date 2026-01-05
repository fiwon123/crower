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

func IsFilePath(path string) bool {

	base := filepath.Base(path)
	splitted := strings.Split(base, ".")

	return len(splitted) == 2
}
