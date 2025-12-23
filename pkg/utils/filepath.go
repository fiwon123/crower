package utils

import (
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
