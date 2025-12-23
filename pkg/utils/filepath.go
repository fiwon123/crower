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

func CreateFolderIfNotExists(path string) {
	if err := os.MkdirAll(path, 0o755); err != nil {
		fmt.Println("Failed to create directory:", err)
	}
}
