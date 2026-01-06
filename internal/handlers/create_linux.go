//go:build linux

package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

func CreateSystemPathVariable(value string, app *app.Data) (string, error) {

	lineSlice := getProfileFileSlice()

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

func getProfileFileSlice() []string {
	home := os.Getenv("HOME")
	profileFile := home + "/.profile"

	lines, err := os.ReadFile(profileFile)
	if err != nil {
		lines = []byte{}
	}

	lineSlice := strings.Split(string(lines), "\n")

	return lineSlice
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
