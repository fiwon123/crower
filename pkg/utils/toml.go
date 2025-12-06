package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/fiwon123/crower/internal/data"
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

func ReadToml(path string) data.CommandsMap {
	dataFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading path file:", err)
		return nil
	}

	var cmds data.CommandsMap
	if err = toml.Unmarshal(dataFile, &cmds); err != nil {
		fmt.Println("Error unmarshal data:", err)
		return nil
	}

	return cmds
}

func WriteToml(cmds data.CommandsMap, path string) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cmds); err != nil {
		fmt.Println("Error enconding path file:", err)
		return
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		fmt.Println("Error writing path file:", err)
		return
	}
}
