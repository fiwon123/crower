// Package sharecommands exports commands to a file, imports them back, and
// reads the shared registry of packs.
//
// Names and aliases share one namespace. Import is the careful part: a pack
// must not shadow a command the user has, so a conflicting name is renamed (or
// replaced by --overwrite), a conflicting alias is renamed or dropped, and a
// collision inside one pack is renamed too. The whole pack is staged before
// anything is applied, and the result is persisted once.
package sharecommands

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
)

// Pack is one entry of the registry index.
type Pack struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

// State is the command set and the way to persist it.
type State interface {
	Commands() catalog.Catalog
	Save(commands catalog.Catalog) error
}

// Output is the diagnostics and the result stream.
type Output interface {
	Info(message string, data ...any)
	Result(text string)
}

// Prompter asks the questions this use case has.
type Prompter interface {
	Confirm(question string) (bool, error)
}

// Packs reads and writes pack files.
type Packs interface {
	Write(path string, commands []catalog.Command) error
	Read(path string) ([]catalog.Command, error)
}

// Registry reads the published packs.
type Registry interface {
	// Index returns every pack the registry lists.
	Index() ([]Pack, error)

	// Pack returns one pack and the URL it came from.
	Pack(name string) ([]catalog.Command, string, error)
}

// UseCase exports, imports and pulls command packs.
type UseCase struct {
	state    State
	output   Output
	prompt   Prompter
	packs    Packs
	registry Registry
}

// New builds the use case from its five dependencies.
func New(state State, output Output, prompt Prompter, packs Packs, registry Registry) *UseCase {
	return &UseCase{state: state, output: output, prompt: prompt, packs: packs, registry: registry}
}

// Export writes one command, or all of them, to a pack file.
func (u *UseCase) Export(name, outPath string) error {
	commands := u.state.Commands()

	if name == "" {
		if outPath == "" {
			outPath = "crower-export.toml"
		}

		return u.packs.Write(outPath, commands.Commands())
	}

	command, found := commands.Resolve(name)
	if !found {
		return fmt.Errorf("command not found")
	}

	if outPath == "" {
		outPath = fmt.Sprintf("%s.toml", name)
	}

	return u.packs.Write(outPath, []catalog.Command{command})
}

// Import reads a pack file and merges it into the command set.
func (u *UseCase) Import(path string, overwrite bool) error {
	pack, err := u.packs.Read(path)
	if err != nil {
		return err
	}

	return u.importPack(pack, overwrite)
}

// List writes the registry's pack lines.
func (u *UseCase) List() error {
	packs, err := u.registry.Index()
	if err != nil {
		return err
	}

	for _, pack := range packs {
		u.output.Result(fmt.Sprintf("%s - %s (v%s)", pack.Name, pack.Description, pack.Version))
	}

	return nil
}

// Pull fetches one pack, shows what it holds, and imports it on a yes.
//
// A preview shown before a confirmation is a diagnostic, not a result: a script
// that pipes the results must not receive the pack body. Only the entry names
// are shown: a pack carries the user's commands, sometimes secrets, and the
// logger writes to a file too.
func (u *UseCase) Pull(name string) error {
	pack, url, err := u.registry.Pack(name)
	if err != nil {
		return err
	}

	names := make([]string, 0, len(pack))
	for _, command := range pack {
		names = append(names, command.Name)
	}

	slices.Sort(names)

	u.output.Info(fmt.Sprintf("Importing pack '%s' from %s", name, url))
	u.output.Info(fmt.Sprintf("Pack contains: %s", strings.Join(names, ", ")))

	proceed, err := u.prompt.Confirm("Import?")
	if err != nil {
		return err
	}

	if !proceed {
		return nil
	}

	return u.importPack(pack, false)
}

// importPack merges the pack into the command set.
//
// Nothing is written before the whole pack is staged: a prompt error or a
// refusal leaves the command set and the file exactly as they were.
func (u *UseCase) importPack(pack []catalog.Command, overwrite bool) error {
	current := u.state.Commands()

	staged := make([]catalog.Command, 0, len(pack))
	stagedNames := map[string]bool{}
	stagedAliases := map[string]bool{}

	slices.SortFunc(pack, func(a, b catalog.Command) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, imported := range pack {
		originalKey := imported.Name
		key := originalKey

		local := map[string]bool{}
		aliases := make([]string, 0, len(imported.Aliases))
		conflicting := make([]bool, 0, len(imported.Aliases))

		nameStaged := stagedNames[key] || stagedAliases[key]
		nameExists := existsIn(current, key)

		// Two entries in one pack still need distinct keys: --overwrite only
		// replaces commands that already exist on disk.
		if overwrite && nameStaged {
			key = firstFree(originalKey, current, stagedNames, stagedAliases, local)
		}

		for _, alias := range imported.Aliases {
			if alias == "" {
				continue
			}

			conflicting = append(conflicting, aliasTaken(alias, key, overwrite, stagedNames, stagedAliases, local, current))
			aliases = append(aliases, alias)
			local[alias] = true
		}

		needRename := (!overwrite && (nameStaged || nameExists)) || slices.Contains(conflicting, true)
		rename := false

		if needRename && !overwrite {
			var err error

			rename, err = u.prompt.Confirm(fmt.Sprintf("Command '%s' conflicts with an existing name or alias. Rename it? (no skips this command)", originalKey))
			if err != nil {
				return err
			}

			if !rename {
				continue
			}
		}

		if rename && (nameStaged || nameExists) {
			key = firstFree(originalKey, current, stagedNames, stagedAliases, local)
		}

		kept := make([]string, 0, len(aliases))

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
			reserved := map[string]bool{key: true}

			renamed := firstFree(alias, current, stagedNames, stagedAliases, local, reserved)
			local[renamed] = true
			kept = append(kept, renamed)
		}

		staged = append(staged, catalog.Command{
			Name:    key,
			Aliases: kept,
			Exec:    imported.Exec,
			Steps:   imported.Steps,
		})

		stagedNames[key] = true

		for _, alias := range kept {
			stagedAliases[alias] = true
		}
	}

	if len(staged) == 0 {
		return nil
	}

	next := current

	for _, command := range staged {
		// --overwrite may be replacing a command with this name, or taking a
		// name that another command used as an alias. Remove that ownership
		// first, or the reaches-old-command alias would survive the import.
		next = takeName(next, command.Name)

		added, err := next.Add(command)
		if err != nil {
			return err
		}

		next = added
	}

	// Persist once, after the whole pack is applied. Import is a bulk setup
	// and records no history entry: one per imported command would describe
	// states that never existed.
	return u.state.Save(next)
}

// existsIn reports whether a name or alias is used anywhere in the set.
func existsIn(commands catalog.Catalog, key string) bool {
	if _, found := commands.Get(key); found {
		return true
	}

	_, found := commands.AliasOwner(key)

	return found
}

// aliasTaken reports whether alias already belongs to a name or an alias
// anywhere: on disk, in an entry staged earlier in this pack, or earlier in the
// same command.
//
// Under --overwrite a key the command being replaced already owns is not a
// conflict: re-importing the same command has to keep its aliases.
func aliasTaken(alias, key string, overwrite bool, stagedNames, stagedAliases, local map[string]bool, current catalog.Catalog) bool {
	if stagedNames[alias] || stagedAliases[alias] || local[alias] {
		return true
	}

	if owner, found := current.AliasOwner(alias); found && !(overwrite && owner == key) {
		return true
	}

	if _, found := current.Get(alias); found && !(overwrite && alias == key) {
		return true
	}

	return false
}

// firstFree is base+"_imported", or the first free numbered variant, so a
// rename never silently replaces a command that already has that name.
func firstFree(base string, current catalog.Catalog, stagedNames, stagedAliases, local map[string]bool, extra ...map[string]bool) string {
	key := base + "_imported"

	for n := 2; taken(key, current, stagedNames, stagedAliases, local, extra...); n++ {
		key = fmt.Sprintf("%s_imported%d", base, n)
	}

	return key
}

func taken(key string, current catalog.Catalog, stagedNames, stagedAliases, local map[string]bool, extra ...map[string]bool) bool {
	if existsIn(current, key) || stagedNames[key] || stagedAliases[key] || local[key] {
		return true
	}

	for _, set := range extra {
		if set[key] {
			return true
		}
	}

	return false
}

// takeName removes any ownership of key so an imported command can take it,
// keeping the set consistent with the records it holds.
func takeName(commands catalog.Catalog, key string) catalog.Catalog {
	if _, found := commands.Get(key); found {
		commands, _, _ = commands.Remove(key)
	}

	return stripAlias(commands, key)
}

// stripAlias removes an alias claim from every command that still lists it.
func stripAlias(commands catalog.Catalog, alias string) catalog.Catalog {
	list := commands.Commands()
	changed := false

	for i := range list {
		kept := slices.DeleteFunc(slices.Clone(list[i].Aliases), func(candidate string) bool {
			return candidate == alias
		})

		if len(kept) != len(list[i].Aliases) {
			changed = true
		}

		list[i].Aliases = kept
	}

	if !changed {
		return commands
	}

	return catalog.New(commands.Names(), list)
}
