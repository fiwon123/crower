package system

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/pkg/utils"
)

// Scripts writes the script file a `create --script` command runs.
//
// The folder is a constructor argument, not a package-level lookup: it is the
// installation's own scripts directory, and this package has no way to know
// where crower keeps its configuration.
type Scripts struct {
	folder string
}

// NewScripts addresses the scripts stored in folder.
func NewScripts(folder string) *Scripts {
	return &Scripts{folder: folder}
}

// Create makes the script file for name, and returns its path and the folder
// that holds it -- which the caller opens so the user can edit the script.
//
// An existing file is kept: creating the command again must not wipe work the
// user put in the script.
func (s *Scripts) Create(name string) (string, string, error) {
	if err := utils.CreateFolderIfNotExists(s.folder); err != nil {
		return "", "", err
	}

	extension, err := ScriptExtension(runtime.GOOS)
	if err != nil {
		return "", "", err
	}

	path := filepath.Join(s.folder, name+extension)

	if err := utils.CreateFileIfNotExists(path); err != nil {
		return "", "", err
	}

	return path, s.folder, nil
}

// ScriptExtension is the suffix of a script command's file: a batch file on
// Windows, a shell script everywhere else.
func ScriptExtension(goos string) (string, error) {
	switch goos {
	case "windows":
		return ".bat", nil
	case "linux", "darwin":
		return ".sh", nil
	default:
		return "", fmt.Errorf("%w: script command on %q", errors.ErrUnsupportedOS, goos)
	}
}
