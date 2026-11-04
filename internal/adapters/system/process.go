package system

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/shell"
	"github.com/fiwon123/crower/pkg/utils"
)

// Processes turns a running process into the line that starts it again.
//
// The lookups are fields rather than direct calls so a test can answer them
// without a process table: a test that had to start a real process would be
// slow, platform-dependent, and would still not cover the flatpak branch.
type Processes struct {
	pathByID   func(int32) (string, error)
	nameByID   func(int32) (string, error)
	pathByName func(string) (string, error)
	flatpakID  func(string) (string, error)
}

// NewProcesses builds a Processes that reads the host's process table.
func NewProcesses() *Processes {
	return &Processes{
		pathByID:   utils.GetProcessPathByID,
		nameByID:   utils.GetProcessNameByID,
		pathByName: utils.GetProcessPathByName,
		flatpakID:  utils.GetFlatpakAppIDByName,
	}
}

// ExecLine resolves a process name or pid to a line a shell can run later.
//
// A flatpak application is started by its app id (`flatpak run <id>`), because
// its executable lives inside a sandbox path that is not the command a user
// would type. Any other executable is used as it is, quoted: the line is stored
// and later handed to a shell.
func (p *Processes) ExecLine(process string) (string, error) {
	processName := ""
	path := ""

	pid, err := strconv.Atoi(process)
	if err != nil {
		processName = process

		path, err = p.pathByName(processName)
		if err != nil {
			return "", err
		}
	} else {
		path, err = p.pathByID(int32(pid))
		if err != nil {
			return "", err
		}
	}

	if strings.Contains(path, "app/") {
		if processName == "" {
			processName, err = p.nameByID(int32(pid))
			if err != nil {
				return "", err
			}
		}

		appID, err := p.flatpakID(processName)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("flatpak run %s", appID), nil
	}

	if path != "" {
		return shell.Quote(path), nil
	}

	return "", fmt.Errorf("couldn't find the process either by pid or by name")
}
