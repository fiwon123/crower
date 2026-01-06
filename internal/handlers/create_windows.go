//go:build windows

package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/fiwon123/crower/internal/data/app"
	"golang.org/x/sys/windows/registry"
)

func CreateSystemPathVariable(value string, app *app.Data) (string, error) {
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

	ok := checkNewVarValue(value, path)
	if !ok {
		return "", fmt.Errorf("value already in PATH")
	}

	newPath := path + string(os.PathListSeparator) + value

	err = key.SetStringValue("Path", newPath)

	return "Added to PATH", err

}

func checkNewVarValue(value string, from string) bool {
	splitted := splitPath(from)
	ok := true
	if slices.Contains(splitted, value) {
		return false
	}

	return ok
}

func splitPath(path string) []string {
	return filepath.SplitList(path)
}
