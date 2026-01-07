//go:build linux

package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

func DeleteSystemVariable(key string, app *app.Data) (string, error) {
	bashrcPath := os.Getenv("HOME") + "/.bashrc"
	fileSlice := getFileLineSlice(bashrcPath)

	found := false
	for i, s := range fileSlice {
		if strings.Contains(s, fmt.Sprintf("export %s=", key)) {
			fileSlice[i] = ""
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("var name not found")
	}

	err := writeFile(fileSlice, bashrcPath)
	if err != nil {
		return "", fmt.Errorf("can't delete var name: %v", err)
	}

	return "var name deleted", nil
}

func DeleteSystemPathVariable(path string, app *app.Data) (string, error) {
	home := os.Getenv("HOME")
	profileFilePath := home + "/.profile"
	lineSlice := getFileLineSlice(profileFilePath)

	pathLinePrefix := "export PATH="
	pathLinePrefix2 := "export PATH"

	found := false

	for i, line := range lineSlice {
		if strings.HasPrefix(line, pathLinePrefix) && strings.Contains(line, path) {
			newLine := strings.Replace(line, path, "", -1)
			lineSlice[i] = cleanPath(newLine, string(os.PathListSeparator))
			found = true
			break
		} else if strings.HasPrefix(line, pathLinePrefix2) {
			break
		}
	}

	if !found {
		return "", fmt.Errorf("path not found")
	}

	err := writeFile(lineSlice, profileFilePath)
	if err != nil {
		return "", err
	}

	return "Deleted path from PATH", err
}

func cleanPath(path string, sep string) string {
	parts := strings.Split(path, sep)
	cleaned := []string{}

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}

	return strings.Join(cleaned, sep)
}
