//go:build linux

package deletedom

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

func (h *Handler) DeleteSystemVariable(key string) (string, error) {
	bashrcPath := os.Getenv("HOME") + "/.bashrc"
	fileSlice := utils.GetFileLineSlice(bashrcPath)

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

	err := utils.WriteFile(fileSlice, bashrcPath)
	if err != nil {
		return "", fmt.Errorf("can't delete var name: %v", err)
	}

	return "var name deleted", nil
}

func (h *Handler) DeleteSystemPathVariable(path string) (string, error) {
	home := os.Getenv("HOME")
	profileFilePath := home + "/.profile"
	lineSlice := utils.GetFileLineSlice(profileFilePath)

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

	err := utils.WriteFile(lineSlice, profileFilePath)
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
