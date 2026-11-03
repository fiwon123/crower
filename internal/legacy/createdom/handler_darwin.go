//go:build darwin

package createdom

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

var pathFiles = []string{".zshrc", ".bashrc", ".profile"}

func (h *Handler) CreateSystemVariable(newVar string, value string) (string, error) {
	zshrcPath := os.Getenv("HOME") + "/.zshrc"
	fileSlice := utils.GetFileLineSlice(zshrcPath)

	for _, s := range fileSlice {
		if strings.Contains(s, fmt.Sprintf("export %s=", newVar)) {
			return "", fmt.Errorf("variable already added")
		}
	}

	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\nexport %s=\"%s\"", newVar, value)
	if err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "Added to .zshrc. Restart terminal to take effect.", nil
}

func (h *Handler) CreateSystemPathVariable(value string) (string, error) {
	home := os.Getenv("HOME")

	added := []string{}
	already := []string{}

	for _, name := range pathFiles {
		changed, err := addPathToFile(home+"/"+name, value)
		if err != nil {
			return "", err
		}

		if changed {
			added = append(added, name)
		} else {
			already = append(already, name)
		}
	}

	if len(added) == 0 {
		return fmt.Sprintf("Already in PATH (%s). Nothing to change.", strings.Join(already, ", ")), nil
	}

	return fmt.Sprintf("Added to PATH in %s. Open a new terminal (or run \"source %s\") for it to take effect.",
		strings.Join(added, ", "), home+"/"+added[0]), nil
}

func addPathToFile(filePath string, value string) (bool, error) {
	lineSlice := utils.GetFileLineSlice(filePath)

	for i, line := range lineSlice {
		if !isPathExport(line) {
			continue
		}

		if pathHasValue(line, value) {
			return false, nil
		}

		lineSlice[i] = addValueToPathExport(line, value)

		err := utils.WriteFile(lineSlice, filePath)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	lineSlice = append(lineSlice, fmt.Sprintf("export PATH=\"$PATH:%s\"", value))

	err := utils.WriteFile(lineSlice, filePath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func isPathExport(line string) bool {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "export PATH=") {
		return true
	}

	return trimmed == "export PATH" || strings.HasPrefix(trimmed, "export PATH ")
}

func pathHasValue(line string, value string) bool {
	for _, element := range pathElements(line) {
		if element == value {
			return true
		}
	}

	return false
}

func pathElements(line string) []string {
	trimmed := strings.TrimSpace(line)

	value := ""
	if index := strings.Index(trimmed, "="); index >= 0 {
		value = trimmed[index+1:]
	}

	value = strings.Trim(value, "\"")

	return strings.Split(value, ":")
}

func addValueToPathExport(line string, value string) string {
	trimmed := strings.TrimSpace(line)

	if !strings.Contains(trimmed, "=") {
		return fmt.Sprintf("export PATH=\"$PATH:%s\"", value)
	}

	start := strings.Index(trimmed, "\"")
	end := strings.LastIndex(trimmed, "\"")

	if start >= 0 && end > start {
		return trimmed[:end] + ":" + value + trimmed[end:]
	}

	return trimmed + ":" + value
}
