package sharedom

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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

	// Stage the whole file before touching the run: a prompt error must leave
	// both the maps and the file as they were. The apply below mutates memory
	// and then persists; the process exits with the import, so only a write
	// failure can leave memory ahead of the file.
	//
	// Names and aliases share one namespace. An import may not shadow either:
	// a conflicting name is renamed (or replaced by --overwrite), and a
	// conflicting alias is renamed, or dropped under --overwrite because
	// taking it would silently edit another command.
	staged := command.NewMapData()
	stagedAliases := command.NewMapData()

	for _, k := range slices.Sorted(maps.Keys(commands)) {
		v := commands[k]
		key := k

		nameStaged := staged.Exists(key) || stagedAliases.Exists(key)
		nameExists := h.app.AllCommandsByName.Exists(key) || h.app.AllCommandsByAlias.Exists(key)

		// Two entries in one pack still need distinct keys: --overwrite only
		// replaces commands that already exist on disk.
		if overwrite && nameStaged {
			key = firstFreeImportKey(k, h.app.AllCommandsByName, h.app.AllCommandsByAlias, staged, stagedAliases)
		}

		local := command.NewMapData()
		aliases := make([]string, 0, len(v.AllAlias))
		conflicting := make([]bool, 0, len(v.AllAlias))

		for _, alias := range v.AllAlias {
			if alias == "" {
				continue
			}

			conflicting = append(conflicting, h.importAliasTaken(alias, key, overwrite, staged, stagedAliases, local))
			aliases = append(aliases, alias)
			local.Add(alias, &command.Data{Name: alias})
		}

		needRename := (!overwrite && (nameStaged || nameExists)) || slices.Contains(conflicting, true)
		rename := false

		if needRename && !overwrite {
			rename, err = helper.GetUserConfirmation(
				fmt.Sprintf("Command '%s' conflicts with an existing name or alias. Rename it? (no skips this command)", k), h.app)
			if err != nil {
				return err
			}

			if !rename {
				continue
			}
		}

		if rename && (nameStaged || nameExists) {
			key = firstFreeImportKey(k, h.app.AllCommandsByName, h.app.AllCommandsByAlias, staged, stagedAliases, local)
		}

		kept := aliases[:0]
		for i, alias := range aliases {
			if !conflicting[i] {
				kept = append(kept, alias)
				continue
			}

			// --overwrite says "do not ask"; an alias that would shadow is
			// dropped rather than taken from the command that owns it.
			if overwrite {
				continue
			}

			// The command's own name is not in the maps yet; reserve it so a
			// renamed alias cannot land on it and create a self-alias.
			reserved := command.NewMapData()
			reserved.Add(key, &command.Data{Name: key})

			renamed := firstFreeImportKey(alias, h.app.AllCommandsByName, h.app.AllCommandsByAlias, staged, stagedAliases, local, reserved)
			local.Add(renamed, &command.Data{Name: renamed})
			kept = append(kept, renamed)
		}

		// The map is keyed by the command's name everywhere else; an imported
		// file must not produce a key and a Name that disagree.
		v.Name = key
		v.AllAlias = kept
		staged.Add(key, &v)

		for _, alias := range kept {
			stagedAliases.Add(alias, &v)
		}
	}

	if len(staged) == 0 {
		return nil
	}

	for _, key := range slices.Sorted(maps.Keys(staged)) {
		cmd := staged.Get(key)

		// --overwrite may be replacing a command with this name, or taking a
		// name that another command used as an alias. Remove that ownership
		// first, or the reaches-old-command alias would survive the import.
		h.takeName(key)

		h.app.AllCommandsByName.Add(key, cmd)

		for _, alias := range cmd.AllAlias {
			h.app.AllCommandsByAlias.Add(alias, cmd)
		}
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

// importAliasTaken reports whether alias already belongs to a name or an alias
// anywhere: on disk, in an entry staged earlier in this pack, or earlier in the
// same command.
//
// Under --overwrite a key the command being replaced already owns is not a
// conflict: re-importing the same command has to keep its aliases.
func (h *Handler) importAliasTaken(alias, key string, overwrite bool, staged, stagedAliases, local command.MapData) bool {
	if staged.Exists(alias) || stagedAliases.Exists(alias) || local.Exists(alias) {
		return true
	}

	if owner := h.app.AllCommandsByAlias.Get(alias); owner != nil && !(overwrite && owner.Name == key) {
		return true
	}

	if owner := h.app.AllCommandsByName.Get(alias); owner != nil && !(overwrite && alias == key) {
		return true
	}

	return false
}

// takeName removes any ownership of key so an imported command can take it,
// keeping both maps consistent with the records they hold.
func (h *Handler) takeName(key string) {
	if old := h.app.AllCommandsByName.Get(key); old != nil {
		h.app.AllCommandsByName.Remove(key)

		for _, alias := range old.AllAlias {
			// Only drop the alias entry when it still points at this command:
			// a legacy duplicate may already have handed it to another one.
			if owner := h.app.AllCommandsByAlias.Get(alias); owner != nil && (owner.Name == old.Name || owner.Name == "") {
				h.app.AllCommandsByAlias.Remove(alias)
			}
		}
	}

	if h.app.AllCommandsByAlias.Get(key) == nil {
		return
	}

	h.app.AllCommandsByAlias.Remove(key)

	// Drop the alias from every record that claims it. Looking the owner up by
	// Name would miss a hand-edited record that omits Name.
	for name, c := range h.app.AllCommandsByName {
		if !slices.Contains(c.AllAlias, key) {
			continue
		}

		c.AllAlias = slices.DeleteFunc(slices.Clone(c.AllAlias), func(alias string) bool {
			return alias == key
		})
		h.app.AllCommandsByName.Update(name, &c)
	}
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
