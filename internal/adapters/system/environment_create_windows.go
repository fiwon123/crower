//go:build windows

package system

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fiwon123/crower/pkg/utils"
	"golang.org/x/sys/windows/registry"
)

// SetVariable stores one value in the user's Environment registry key.
func (e *Environment) SetVariable(name, value string) (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Environment`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return "", fmt.Errorf("open environment key: %w", err)
	}
	defer key.Close()

	if _, _, err := key.GetStringValue(name); err == nil {
		return "", fmt.Errorf("value already exists")
	}

	cmd := exec.Command("reg", "add", `HKCU\Environment`, "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Create System Variable error: %v \n", err)
	}

	return "System variable set. You may need to restart or log off for it to take effect.", nil
}

// AddPath appends one directory to the user's Path value.
func (e *Environment) AddPath(value string) (string, error) {
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

	if !utils.CheckNewVarValuePath(value, path) {
		return "Already in PATH. Nothing to change.", nil
	}

	newPath := path + string(os.PathListSeparator) + value

	if err := key.SetStringValue("Path", newPath); err != nil {
		return "", err
	}

	return "Added to PATH. Open a new terminal (or sign out and in) for it to take effect.", nil
}
