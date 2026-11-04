//go:build darwin

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

// pathFiles are the files the PATH export is added to: zsh is the login shell
// on current macOS, and the bash and POSIX files are kept for setups that still
// read them.
var pathFiles = []string{".zshrc", ".bashrc", ".profile"}

// SetVariable appends one `export NAME="value"` line to ~/.zshrc.
func (e *Environment) SetVariable(name, value string) (string, error) {
	zshrcPath := os.Getenv("HOME") + "/.zshrc"

	for _, line := range utils.GetFileLineSlice(zshrcPath) {
		if strings.Contains(line, fmt.Sprintf("export %s=", name)) {
			return "", fmt.Errorf("variable already added")
		}
	}

	// The file may not exist yet on a fresh account. The legacy code called
	// log.Fatal here, killing the process; appending to a file the operation
	// is about to fill is the obvious intent, so create it.
	if err := utils.CreateFileIfNotExists(zshrcPath); err != nil {
		return "", err
	}

	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", zshrcPath, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "\nexport %s=\"%s\"", name, value); err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "Added to .zshrc. Restart terminal to take effect.", nil
}

// AddPath adds a directory to the PATH export of every file the shell reads.
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
