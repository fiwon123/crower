//go:build linux

package createdom

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

// pathFiles are the files the PATH export is added to. A login shell reads
// .profile, an interactive shell reads .bashrc, and which one a system reads
// is not the same everywhere: .profile is a Debian/Ubuntu convention, and on
// Fedora and Arch the login file is .bash_profile, which sources .bashrc.
var pathFiles = []string{".profile", ".bashrc"}

func (h *Handler) CreateSystemVariable(newVar string, value string) (string, error) {

	bashrcPath := os.Getenv("HOME") + "/.bashrc"
	fileSlice := utils.GetFileLineSlice(bashrcPath)

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

// CreateSystemPathVariable adds a directory to the PATH of the user. It
// reports where the export was written and what has to happen for the
// current shell to see it.
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

// addPathToFile adds the value to the PATH export of the file. It reports
// whether the file was changed: a file that already has the value in its
// PATH is left alone.
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

// isPathExport reports whether the line is a PATH export. A line that starts
// with "export PATH" but is not one — "export PATHDIR=…" — is not.
func isPathExport(line string) bool {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "export PATH=") {
		return true
	}

	return trimmed == "export PATH" || strings.HasPrefix(trimmed, "export PATH ")
}

// pathHasValue reports whether the PATH export already has the value as one
// of its elements. The elements are compared whole: "/opt/tools" is not in a
// PATH that holds "/opt/tools-extra".
func pathHasValue(line string, value string) bool {
	for _, element := range pathElements(line) {
		if element == value {
			return true
		}
	}

	return false
}

// pathElements returns the elements of the PATH export of the line.
func pathElements(line string) []string {

	trimmed := strings.TrimSpace(line)

	value := ""
	if index := strings.Index(trimmed, "="); index >= 0 {
		value = trimmed[index+1:]
	}

	value = strings.Trim(value, "\"")

	return strings.Split(value, ":")
}

// addValueToPathExport returns the PATH export with the value added to it,
// keeping the form the line already uses.
func addValueToPathExport(line string, value string) string {

	trimmed := strings.TrimSpace(line)

	// A line that is only "export PATH" carries no value to extend.
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
