//go:build linux

package createdom

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/crowerutils"
)

func (h *CreateHandler) CreateSystemVariable(newVar string, value string) (string, error) {

	bashrcPath := os.Getenv("HOME") + "/.bashrc"
	fileSlice := crowerutils.GetFileLineSlice(bashrcPath)

	for _, s := range fileSlice {
		if strings.Contains(s, fmt.Sprintf("export %s=", newVar)) {
			return "", fmt.Errorf("variable already added")
		}
	}

	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\nexport %s=\"%s\"", newVar, value)
	if err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "Added to .bashrc. Restart terminal to take effect.", nil
}

func (h *CreateHandler) CreateSystemPathVariable(value string) (string, error) {

	home := os.Getenv("HOME")
	profileFilePath := home + "/.profile"
	lineSlice := crowerutils.GetFileLineSlice(profileFilePath)

	pathLinePrefix := "export PATH="
	pathLinePrefix2 := "export PATH"

	found := false

	for i, line := range lineSlice {
		if strings.HasPrefix(line, pathLinePrefix) {
			if !strings.Contains(line, value) {
				start := strings.Index(line, "\"")
				end := strings.LastIndex(line, "\"")

				if start >= 0 && end > start {
					lineSlice[i] = line[:end] + ":" + value + line[end:]
				} else {
					lineSlice[i] = line + ":" + value
				}
			}
			found = true
			break
		} else if strings.HasPrefix(line, pathLinePrefix2) {
			if !strings.Contains(line, value) {
				lineSlice[i] = line[:] + "=\"$PATH:" + value + "\""
			}
			found = true
			break
		}
	}

	if !found {
		lineSlice = append(lineSlice, fmt.Sprintf("export PATH=\"$PATH:%s\"", value))
	}

	err := crowerutils.WriteFile(lineSlice, profileFilePath)
	if err != nil {
		return "", err
	}

	return "Added to PATH", err
}
