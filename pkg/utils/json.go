package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func ReadJson(filePath string, output any) error {
	dataFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file path: %v", err)
	}

	if len(bytes.TrimSpace(dataFile)) == 0 {
		return nil
	}

	err = json.Unmarshal(dataFile, output)
	if err != nil {
		return fmt.Errorf("Error unmarshal json file: %v", err)
	}

	return nil
}

func WriteJson(input string, filePath string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(input); err != nil {
		return fmt.Errorf("Error enconding file path: %v", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Error writing json file: %v", err)
	}

	return nil
}
