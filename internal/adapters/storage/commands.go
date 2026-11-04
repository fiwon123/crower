package storage

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/adapters/storage/command"

	"github.com/fiwon123/crower/pkg/utils"
)

// CommandFile is the TOML file that holds the managed commands.
type CommandFile struct {
	path string
}

// NewCommandFile addresses the command map at path.
func NewCommandFile(path string) *CommandFile {
	return &CommandFile{path: path}
}

// Path is where the command map lives.
func (f *CommandFile) Path() string {
	return f.path
}

// Folder is the directory the configuration lives in, which is also where
// `create --script` keeps its scripts.
func (f *CommandFile) Folder() string {
	return filepath.Dir(f.path)
}

// Load reads the commands and the order their keys appear in.
//
// It prepares the folder and the file first, so a first run starts from an
// empty set instead of an error. Both reads are checked: a config whose key
// order could not be read must not lose the order of every command without a
// word.
func (f *CommandFile) Load() (command.MapData, []string, error) {
	commands := command.NewMapData()

	if err := utils.CreateFolderIfNotExists(filepath.Dir(f.path)); err != nil {
		return nil, nil, fmt.Errorf("prepare config folder: %w", err)
	}

	if err := utils.CreateFileIfNotExists(f.path); err != nil {
		return nil, nil, fmt.Errorf("prepare config file %s: %w", f.path, err)
	}

	order, err := utils.ReadKeysTomlInOrder(f.path)
	if err != nil {
		return nil, nil, fmt.Errorf("read command order from %s: %w", f.path, err)
	}

	if err := utils.ReadToml(f.path, &commands); err != nil {
		return nil, nil, fmt.Errorf("read config %s: %w", f.path, err)
	}

	return commands, order, nil
}

// Save writes the whole command map, in display order.
func (f *CommandFile) Save(commands command.MapData) error {
	return utils.WriteToml(commands, f.path)
}
