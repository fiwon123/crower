package sharedom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/pelletier/go-toml"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{app: app}
}

type RegistryIndex struct {
	Packs []RegistryPack `json:"packs"`
}

type RegistryPack struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

func (h *Handler) Export(name string, outPath string) error {
	if name == "" {
		// Export all
		data := h.app.AllCommandsByName
		out, err := toml.Marshal(data)
		if err != nil {
			return err
		}
		if outPath == "" {
			outPath = "crower-export.toml"
		}
		return os.WriteFile(outPath, out, 0644)
	}

	cmd := h.app.AllCommandsByName.Get(name)
	if cmd == nil {
		cmd = h.app.AllCommandsByAlias.Get(name)
	}
	if cmd == nil {
		return fmt.Errorf("command not found")
	}

	data := map[string]*command.Data{cmd.Name: cmd}
	out, err := toml.Marshal(data)
	if err != nil {
		return err
	}
	if outPath == "" {
		outPath = fmt.Sprintf("%s.toml", name)
	}
	return os.WriteFile(outPath, out, 0644)
}

func (h *Handler) Import(inPath string, overwrite bool) error {
	data, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	var commands map[string]command.Data
	if err := toml.Unmarshal(data, &commands); err != nil {
		return err
	}

	// Stage the whole file before touching the run: a prompt or a write that
	// fails must leave both the map and the file as they were.
	staged := command.NewMapData()

	for k, v := range commands {
		key := k

		if importKeyTaken(key, h.app.AllCommandsByName, staged) && !overwrite {
			rename, err := helper.GetUserConfirmation(fmt.Sprintf("Command '%s' exists. Rename?", k), h.app)
			if err != nil {
				return err
			}

			if !rename {
				continue
			}

			key = firstFreeImportKey(k, h.app.AllCommandsByName, staged)
		}

		// The map is keyed by the command's name everywhere else; an imported
		// file must not produce a key and a Name that disagree.
		v.Name = key
		staged.Add(key, &v)
	}

	if len(staged) == 0 {
		return nil
	}

	for key, cmd := range staged {
		h.app.AllCommandsByName.Add(key, &cmd)
	}

	// Persist once, after the whole file is applied. The old code called Save
	// per command, which writes a history snapshot and the index -- never the
	// command map -- so the imported commands vanished on the next run, and it
	// overwrote the snapshot of the operation before the import.
	//
	// No history entry is recorded: import is a bulk setup, and one entry per
	// imported command would describe states that never existed.
	return h.app.SaveCommands()
}

// importKeyTaken reports whether key is used in any of the maps.
func importKeyTaken(key string, maps ...command.MapData) bool {
	for _, m := range maps {
		if m.Exists(key) {
			return true
		}
	}

	return false
}

// firstFreeImportKey is base+"_imported", or the first free numbered variant,
// so a rename never silently replaces a command that already has that name.
func firstFreeImportKey(base string, maps ...command.MapData) string {
	key := base + "_imported"

	for n := 2; importKeyTaken(key, maps...); n++ {
		key = fmt.Sprintf("%s_imported%d", base, n)
	}

	return key
}

func (h *Handler) getRegistryURL() string {
	return "https://raw.githubusercontent.com/fiwon123/crower-resources/main"
}

func (h *Handler) ListRegistry() error {
	url := h.getRegistryURL() + "/index.json"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch registry: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var idx RegistryIndex
	if err := json.Unmarshal(body, &idx); err != nil {
		return err
	}

	for _, p := range idx.Packs {
		fmt.Printf("%s - %s (v%s)\n", p.Name, p.Description, p.Version)
	}
	return nil
}

func (h *Handler) Pull(name string) error {
	if !strings.HasPrefix(h.getRegistryURL(), "https://") {
		return fmt.Errorf("only HTTPS allowed for security")
	}
	url := h.getRegistryURL() + "/packs/" + name + ".toml"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("pack not found: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Save to temp file and import
	tmp := filepath.Join(os.TempDir(), name+".toml")
	if err := os.WriteFile(tmp, body, 0644); err != nil {
		return err
	}
	defer os.Remove(tmp)

	// A preview shown before a confirmation is a diagnostic, not a result: a
	// script that pipes the results must not receive the pack body.
	h.app.Logger.Info(fmt.Sprintf("Importing pack '%s' from %s", name, url))
	h.app.Logger.Info(fmt.Sprintf("Contents:\n%s", string(body)))

	ok, err := helper.GetUserConfirmation("Import?", h.app)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	return h.Import(tmp, false)
}
