package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"
	tomlu "github.com/pelletier/go-toml/v2/unstable"
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

func ReadKeysTomlInOrder(path string) ([]string, error) {
	dataFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading path file: %v", err)
	}

	var keysInOrder []string

	p := tomlu.Parser{}
	p.Reset(dataFile)

	for p.NextExpression() {
		e := p.Expression()

		if e.Kind != tomlu.Table {
			continue
		}

		it := e.Key()
		parts := keyAsStrings(it)

		keysInOrder = append(keysInOrder, string(parts[0]))
	}

	return keysInOrder, nil
}

func keyAsStrings(it tomlu.Iterator) []string {
	var parts []string
	for it.Next() {
		n := it.Node()
		parts = append(parts, string(n.Data))
	}
	return parts
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
