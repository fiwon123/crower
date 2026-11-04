//go:build linux

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

// Variable removes one `export NAME=` line from ~/.bashrc.
//
// The line is emptied rather than removed, which is what `delete --system` has
// always done: it keeps the rest of the file's line numbering intact and leaves
// any continuation line where the user put it.
func (e *Environment) Variable(key string) (string, error) {
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

	if err := utils.WriteFile(fileSlice, bashrcPath); err != nil {
		return "", fmt.Errorf("can't delete var name: %v", err)
	}

	return "var name deleted", nil
}

// VariablePath removes one entry from the export PATH= line in ~/.profile.
func (e *Environment) VariablePath(path string) (string, error) {
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

	if err := utils.WriteFile(lineSlice, profileFilePath); err != nil {
		return "", err
	}

	return "Deleted path from PATH", nil
}

// cleanPath drops the empty segments a removed entry leaves behind.
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
