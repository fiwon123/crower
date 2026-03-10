package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	tomlu "github.com/pelletier/go-toml/v2/unstable"
)

func ReadToml(filePath string, output any) error {
	dataFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file path: %v", err)
	}

	if err = toml.Unmarshal(dataFile, output); err != nil {
		return fmt.Errorf("Error unmarshal data: %v", err)
	}

	return nil
}

func ReadKeysTomlInOrder(filePath string) ([]string, error) {
	dataFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file path: %v", err)
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

func WriteToml(input any, filePath string) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(input); err != nil {
		return fmt.Errorf("Error enconding file path: %v", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Error writing toml file: %v", err)
	}

	return nil
}
