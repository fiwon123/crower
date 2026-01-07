//go:build windows

package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/fiwon123/crower/internal/data/app"
	"golang.org/x/sys/windows/registry"
)

func CreateSystemVariable(newVar string, value string, app *app.Data) (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Environment`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		panic(err)
	}
	defer key.Close()

	_, _, err = key.GetStringValue(newVar)
	if err == nil {
		return "", fmt.Errorf("value already exists")
	}

	cmd := exec.Command("reg", "add", `HKCU\Environment`, "/v", newVar, "/t", "REG_SZ", "/d", value, "/f")
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "System variable set. You may need to restart or log off for it to take effect.", nil

}

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
