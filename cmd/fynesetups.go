package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/session"
)

// defaultSetupName is what the original installation is called in the window.
const defaultSetupName = "default"

// fyneSetups lists and switches the named setups for the window.
//
// It has no Fyne import, so the adapter is covered by the normal cmd tests.
type fyneSetups struct {
	cfg     *session.Session
	current string
}

func newFyneSetups(cfg *session.Session) fyneSetups {
	return fyneSetups{cfg: cfg, current: currentSetupName()}
}

// Current is the setup this window opened with.
func (s fyneSetups) Current() string {
	return s.current
}

// Names lists every setup, the default one first.
func (s fyneSetups) Names() []string {
	return setupNames()
}

// Switch makes name the default for future runs and reloads the session.
//
// The session value is reloaded in place, so every port the window holds keeps
// working against the new setup.
func (s fyneSetups) Switch(name string) (desktop.State, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = defaultSetupName
	}

	path, err := resolveConfigPath("", name)
	if err != nil {
		return desktop.State{}, err
	}

	if err := persistSetup(name); err != nil {
		return desktop.State{}, err
	}

	if err := s.cfg.Reload(path); err != nil {
		return desktop.State{}, err
	}

	return desktop.State{Catalog: s.cfg.Catalog(), History: s.cfg.History}, nil
}

// Commands counts the commands a setup holds.
//
// It checks for the file first: the command loader prepares a missing file so a
// first run starts clean, which would create a setup merely by looking at it.
func (s fyneSetups) Commands(name string) (int, error) {
	path, err := resolveConfigPath("", name)
	if err != nil {
		return 0, err
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}

		return 0, err
	}

	file, _, err := storage.NewCommandFile(path).Load()
	if err != nil {
		return 0, err
	}

	return len(file), nil
}

// Create makes an empty named setup, the way `crower setup new` does.
func (s fyneSetups) Create(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	_, err = createSetup(home, strings.TrimSpace(name))

	return err
}

// Delete removes a named setup and its folder, refusing the default and the
// active setup just as the CLI does without --force.
func (s fyneSetups) Delete(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return deleteSetup(home, strings.TrimSpace(name), s.current)
}

// setupNames lists the original installation and the folders the setups live
// in. The default entry comes first and the rest are sorted, so the list does
// not depend on the order the filesystem returns.
func setupNames() []string {
	names := []string{defaultSetupName}

	home, err := os.UserHomeDir()
	if err != nil {
		return names
	}

	entries, err := os.ReadDir(filepath.Join(home, "crower", "setups"))
	if err != nil {
		return names
	}

	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	slices.Sort(names[1:])

	return names
}

// persistSetup writes the marker `crower setup use` writes, so the next CLI run
// opens the same setup the window chose.
func persistSetup(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	marker := filepath.Join(home, "crower", "setup")

	if err := os.MkdirAll(filepath.Dir(marker), 0o755); err != nil {
		return err
	}

	return os.WriteFile(marker, []byte(name+"\n"), 0o644)
}

// currentSetupName is the setup this process is using, with the window's name
// for the unnamed one.
func currentSetupName() string {
	if name := strings.TrimSpace(activeSetup()); name != "" {
		return name
	}

	return defaultSetupName
}
