package storage

import (
	"fmt"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/pkg/utils"
)

// Packs reads and writes command packs: the TOML files `export` produces and
// `import` consumes, keyed by command name.
type Packs struct{}

// NewPacks builds a Packs.
func NewPacks() *Packs {
	return &Packs{}
}

// Write stores commands at path.
//
// The file is keyed by the command's name, which is the shape the command file
// itself uses and what an older crower reads.
func (p *Packs) Write(path string, commands []catalog.Command) error {
	pack := make(map[string]*command.Data, len(commands))

	for _, item := range commands {
		pack[item.Name] = &command.Data{
			Name:        item.Name,
			AllAlias:    item.Aliases,
			Exec:        item.Exec,
			Steps:       item.Steps,
			Description: item.Description,
			Params:      item.Params,
		}
	}

	if err := utils.WriteToml(pack, path); err != nil {
		return fmt.Errorf("write pack %s: %w", path, err)
	}

	return nil
}

// Read loads the pack at path.
//
// The map key wins over the stored Name, because that is the name every other
// part of crower uses for the command -- a file whose key and Name disagree
// would otherwise import a command under one name and index it under another.
func (p *Packs) Read(path string) ([]catalog.Command, error) {
	pack := map[string]command.Data{}

	if err := utils.ReadToml(path, &pack); err != nil {
		return nil, fmt.Errorf("read pack %s: %w", path, err)
	}

	return fromPack(pack), nil
}

// fromPack converts a decoded pack to the domain shape, keyed by the file key.
func fromPack(pack map[string]command.Data) []catalog.Command {
	commands := make([]catalog.Command, 0, len(pack))

	for key, value := range pack {
		commands = append(commands, catalog.Command{
			Name:        key,
			Aliases:     value.AllAlias,
			Exec:        value.Exec,
			Steps:       value.Steps,
			Description: value.Description,
			Params:      value.Params,
		})
	}

	return commands
}
