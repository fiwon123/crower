package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Function to get Flatpak app ID by name
func GetFlatpakAppIDByName(appName string) (string, error) {
	cmd := exec.Command("flatpak", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Iterate through the lines to find the app's App ID
	lines := bytes.SplitSeq(out, []byte("\n"))
	for line := range lines {
		if bytes.Contains(line, []byte(appName)) {
			// Extract the App ID from the first column
			fields := strings.Fields(string(line))
			if len(fields) > 0 {
				return fields[1], nil
			}
		}
	}
	return "", fmt.Errorf("App '%s' not found", appName)
}
