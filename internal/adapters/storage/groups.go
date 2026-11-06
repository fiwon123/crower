package storage

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/fiwon123/crower/pkg/utils"
)

// groupsName is the file that holds a group's icon, beside the configuration.
const groupsName = "groups.toml"

// groupIcon is one entry of that file.
type groupIcon struct {
	Name string
	Icon string
}

// groupsFile is the on-disk shape: a list of entries rather than a map.
//
// A flat map would encode "v1.2" as a nested table, which the reader could not
// decode back into a string -- a group name with a dot would quietly corrupt
// the file. A list of tables quotes the name and keeps the shape.
type groupsFile struct {
	Groups []groupIcon `toml:"group"`
}

// Groups stores the groups the window offers and the icon each shows.
//
// It is a second file rather than a table inside the command file: the loader
// treats every top-level table as a command, and a mixed shape would fail the
// command-map decode. The CLI never writes this file; it is the window's
// presentation, and a missing file simply means "there are no stored groups".
//
// A group can exist before a command uses it, which is what lets the window
// create one from its own dialog: the entry is kept even when it carries no
// icon, so the name survives a restart that no command refers to.
type Groups struct {
	path string
}

// NewGroups addresses the group icons of the installation stored in folder.
func NewGroups(folder string) *Groups {
	return &Groups{path: filepath.Join(folder, groupsName)}
}

// Path is where the group icons live.
func (g *Groups) Path() string {
	return g.path
}

// Load reads the groups and their icons. A missing file is an empty map, not an
// error: a fresh installation has nothing to say about groups.
func (g *Groups) Load() (map[string]string, error) {
	icons := make(map[string]string)

	if _, err := os.Stat(g.path); errors.Is(err, os.ErrNotExist) {
		return icons, nil
	}

	var stored groupsFile

	if err := utils.ReadToml(g.path, &stored); err != nil {
		return nil, fmt.Errorf("read group icons %s: %w", g.path, err)
	}

	for _, entry := range stored.Groups {
		if entry.Name == "" {
			continue
		}

		icons[entry.Name] = entry.Icon
	}

	return icons, nil
}

// Save writes the groups, sorted by name so the file is stable between runs.
// A group with no icon is still written: the group is the record, the icon is
// decoration.
func (g *Groups) Save(icons map[string]string) error {
	stored := groupsFile{Groups: make([]groupIcon, 0, len(icons))}

	for _, name := range slices.Sorted(maps.Keys(icons)) {
		if name == "" {
			continue
		}

		stored.Groups = append(stored.Groups, groupIcon{Name: name, Icon: icons[name]})
	}

	if err := utils.WriteToml(stored, g.path); err != nil {
		return fmt.Errorf("write group icons %s: %w", g.path, err)
	}

	return nil
}

// Set stores one group and its icon and returns the whole map. It creates the
// group when it is new; clearing the icon keeps the group.
func (g *Groups) Set(group, icon string) (map[string]string, error) {
	icons, err := g.Load()
	if err != nil {
		return nil, err
	}

	if group == "" {
		return icons, nil
	}

	icons[group] = icon

	if err := g.Save(icons); err != nil {
		return nil, err
	}

	return icons, nil
}

// Remove deletes one group and returns the whole map. A group the commands
// still use comes back without an icon, because the catalog derives it.
func (g *Groups) Remove(group string) (map[string]string, error) {
	icons, err := g.Load()
	if err != nil {
		return nil, err
	}

	delete(icons, group)

	if err := g.Save(icons); err != nil {
		return nil, err
	}

	return icons, nil
}
