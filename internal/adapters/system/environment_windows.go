//go:build windows

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
	"golang.org/x/sys/windows/registry"
)

// Variable removes one value from the user's Environment registry key.
func (e *Environment) Variable(varName string) (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Environment`,
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		return "", fmt.Errorf("open environment key: %w", err)
	}
	defer key.Close()

	if _, _, err := key.GetStringValue(varName); err != nil {
		return "", fmt.Errorf("var name not found")
	}

	if err := key.DeleteValue(varName); err != nil {
		return "", fmt.Errorf("can't delete var name: %v", err)
	}

	return "var name deleted", nil
}

// VariablePath removes one entry from the user's Path value.
func (e *Environment) VariablePath(pathValue string) (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return "", err
	}
	defer key.Close()

	path, _, err := key.GetStringValue("Path")
	if err != nil {
		return "", err
	}

	if utils.CheckNewVarValuePath(pathValue, path) {
		return "", fmt.Errorf("path not found in PATH")
	}

	pathSlice := utils.SplitPath(path)
	pathSlice = removeFromSlice(pathSlice, pathValue)

	newPath := strings.Join(pathSlice, string(os.PathListSeparator))

	if err := key.SetStringValue("Path", newPath); err != nil {
		return "", fmt.Errorf("can't delete path from PATH: %v", err)
	}

	return "Deleted path from PATH", nil
}

// removeFromSlice drops every occurrence of item and the empty entries a
// trailing separator leaves behind.
func removeFromSlice(slice []string, item string) []string {
	newSlice := []string{}

	for _, s := range slice {
		if s != item && s != "" {
			newSlice = append(newSlice, s)
		}
	}

	return newSlice
}
