//go:build linux

package handlers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

func CreateSystemVariable(newVar string, value string, app *app.Data) (string, error) {

	bashrcPath := os.Getenv("HOME") + "/.bashrc"
	fileSlice := getFileSlice(bashrcPath)

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

	_, err = fmt.Fprintf(f, "\nexport %s=%s", newVar, value)
	if err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "Added to .bashrc. Restart terminal to take effect.", nil
}

func CreateSystemPathVariable(value string, app *app.Data) (string, error) {

	home := os.Getenv("HOME")
	profileFile := home + "/.profile"
	lineSlice := getFileSlice(profileFile)

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
				lineSlice[i] = line[:len(line)] + "=\"$PATH:" + value + "\""
			}
			found = true
			break
		}
	}

	if !found {
		lineSlice = append(lineSlice, fmt.Sprintf("export PATH=\"$PATH:%s\"", value))
	}

	err := writeProfileFile(lineSlice)
	if err != nil {
		return "", err
	}

	return "Added to PATH", err
}

func writeProfileFile(lineSlice []string) error {
	home := os.Getenv("HOME")
	profileFile := home + "/.profile"

	err := os.WriteFile(profileFile, []byte(strings.Join(lineSlice, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("Error writing .profile: %v", err)
	}

	return nil
}

func lineExists(file, line string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == line {
			return true
		}
	}
	return false
}
