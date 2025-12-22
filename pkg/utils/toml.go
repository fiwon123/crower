package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

func CreateTomlIfNotExists(path string) {
	file, err := os.OpenFile(path,
		os.O_CREATE|os.O_RDWR, // Create if not exists, read/write
		0644)                  // File permissions
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func ReadToml(path string, output any) error {
	dataFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading path file: %v", err)
	}

	if err = toml.Unmarshal(dataFile, output); err != nil {
		return fmt.Errorf("Error unmarshal data: %v", err)
	}

	return nil
}

func WriteToml(input any, path string) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(input); err != nil {
		return fmt.Errorf("Error enconding path file: %v", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Error writing path file: %v", err)
	}

	return nil
}
