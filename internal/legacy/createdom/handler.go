package createdom

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/errors"
)

// Handler owns the rules that decide whether a command may be added. The
// platform work -- creating paths, writing environment files, resolving a
// process -- lives in internal/adapters/system.
type Handler struct {
	app         *app.Config
	creator     *system.Creator
	environment *system.Environment
	processes   *system.Processes
	scripts     *system.Scripts
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{
		app:         app,
		creator:     system.NewCreator(),
		environment: system.NewEnvironment(),
		processes:   system.NewProcesses(),
		scripts:     system.NewScripts(filepath.Join(app.ConfigFolder(), "scripts")),
	}
}

// Create command using name, alias and exec parameters
func (h *Handler) CreateCommand(name string, alias []string, exec string) (*command.Data, error) {

	command := command.New(name, alias, exec)

	if command.Name == "" {
		return nil, fmt.Errorf("empty name")
	}

	if command.Exec == "" {
		return nil, fmt.Errorf("empty exec")
	}

	if h.app.AllCommandsByName.Get(command.Name) != nil {
		return nil, fmt.Errorf("found name, command already added")
	}

	// A name that another command already answers to as an alias would shadow
	// it: lookups try the name map first, so the command the alias pointed at
	// becomes unreachable. The alias loop below cannot catch this one, because
	// it checks the new command's aliases, not its name.
	if owner := h.app.AllCommandsByAlias.Get(command.Name); owner != nil {
		return nil, fmt.Errorf("%q is already an alias of %q", command.Name, owner.Name)
	}

	for _, alias := range command.AllAlias {
		if h.app.AllCommandsByAlias.Get(alias) != nil || h.app.AllCommandsByName.Get(alias) != nil {
			return nil, fmt.Errorf("found alias, command already added")
		}
	}

	h.app.AllCommandsByName.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		h.app.AllCommandsByAlias.Add(alias, command)
	}

	return command, nil
}

// Create command based on process name or id process.
//
// name is the command to create; args[0] is the process name or pid. The
// missing-process case used to index an empty slice and panic, so `crower
// create -p <name>` crashed instead of saying what was missing.
func (h *Handler) CreateProcess(name string, allAlias []string, args []string) (*command.Data, error) {
	if len(args) > 0 && name == "" {
		name = args[0]
		args = args[1:]
	}

	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return nil, errors.NotEnoughArgs("process name or pid")
	}

	execLine, err := h.processes.ExecLine(args[0])
	if err != nil {
		return nil, err
	}

	return h.CreateCommand(name, allAlias, execLine)
}

// Create a new file on filepath
func (h *Handler) CreateFile(filePath string) error {
	return h.creator.File(filePath)
}

// Create a new folder on folderpath
func (h *Handler) CreateFolder(folderPath string) error {
	return h.creator.Folder(folderPath)
}

// CreateScriptCommand writes the script file for name and returns its path.
func (h *Handler) CreateScriptCommand(name string) (string, error) {
	path, _, err := h.scripts.Create(name)
	return path, err
}
