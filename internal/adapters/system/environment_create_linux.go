//go:build linux

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

// pathFiles are the files the PATH export is added to. A login shell reads
// .profile, an interactive shell reads .bashrc, and which one a system reads
// is not the same everywhere: .profile is a Debian/Ubuntu convention, and on
// Fedora and Arch the login file is .bash_profile, which sources .bashrc.
var pathFiles = []string{".profile", ".bashrc"}

// SetVariable appends one `export NAME="value"` line to ~/.bashrc.
func (e *Environment) SetVariable(name, value string) (string, error) {
	bashrcPath := os.Getenv("HOME") + "/.bashrc"

	for _, line := range utils.GetFileLineSlice(bashrcPath) {
		if strings.Contains(line, fmt.Sprintf("export %s=", name)) {
			return "", fmt.Errorf("variable already added")
		}
	}

	// The file may not exist yet on a fresh account. The legacy code called
	// log.Fatal here, killing the process; appending to a file the operation
	// is about to fill is the obvious intent, so create it.
	if err := utils.CreateFileIfNotExists(bashrcPath); err != nil {
		return "", err
	}

	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", bashrcPath, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "\nexport %s=\"%s\"", name, value); err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "Added to .bashrc. Restart terminal to take effect.", nil
}

// AddPath adds a directory to the PATH export of every file the shell reads.
// It reports where the export was written and what has to happen for the
// current shell to see it.
func (e *Environment) AddPath(value string) (string, error) {
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
