//go:build windows

package deletedom

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
	"golang.org/x/sys/windows/registry"
)

func (h *Handler) DeleteSystemVariable(varName string) (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Environment`,
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		panic(err)
	}
	defer key.Close()

	_, _, err = key.GetStringValue(varName)
	if err != nil {
		return "", fmt.Errorf("var name not found")
	}

	err = key.DeleteValue(varName)
	if err != nil {
		return "", fmt.Errorf("can't delete var name: %v", err)
	}

	return "var name deleted", nil
}

func (h *Handler) DeleteSystemPathVariable(pathValue string) (string, error) {
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

	ok := utils.CheckNewVarValuePath(pathValue, path)
	if ok {
		return "", fmt.Errorf("path not found in PATH")
	}

	pathSlice := utils.SplitPath(path)
	pathSlice = removeFromSlice(pathSlice, pathValue)

	newPath := strings.Join(pathSlice, string(os.PathListSeparator))
	err = key.SetStringValue("Path", newPath)
	if err != nil {
		return "", fmt.Errorf("can't delete path from PATH: %v", err)
	}

	return "Deleted path from PATH", nil
}

func removeFromSlice(slice []string, item string) []string {
	newSlice := []string{}
	for _, s := range slice {
		if s != item && s != "" {
			newSlice = append(newSlice, s)
		}
	}
	return newSlice
}
