package sharedom_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/sharedom"
)

const greetPack = `[greet]
Name = "greet"
AllAlias = ["hi"]
Exec = "echo hi"
`

func newConfig(t *testing.T, stdin string) (*app.Config, string) {
	t.Helper()

	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	cfg, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(stdin), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	return cfg, cfgPath
}

func reload(t *testing.T, cfgPath string) *app.Config {
	t.Helper()

	cfg, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("reload: %v", err)
	}

	return cfg
}

func writePack(t *testing.T, contents string) string {
	t.Helper()

	pack := filepath.Join(t.TempDir(), "pack.toml")
	if err := os.WriteFile(pack, []byte(contents), 0o644); err != nil {
		t.Fatalf("write pack: %v", err)
	}

	return pack
}

// Import has to write the configuration file: Save writes a history snapshot
// and the index, not the command map, so the commands used to vanish on the
// next run while the command exited 0.
func TestImportWritesTheConfigurationFile(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	if err := sharedom.NewHandler(cfg).Import(writePack(t, greetPack), false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	// A second run reads the file, the way the next process would.
	imported := reload(t, cfgPath).AllCommandsByName.Get("greet")
	if imported == nil {
		t.Fatal("the imported command was not written to the configuration file")
	}

	if imported.Exec != "echo hi" {
		t.Errorf("exec = %q, want %q", imported.Exec, "echo hi")
	}

	if !slices.Contains(reload(t, cfgPath).OrderKeys, "greet") {
		t.Error("the imported command is not in the file's key order")
	}
}

// Renaming has to insert the command under the new key, not overwrite the key
// it collided with while carrying a different Name.
func TestImportRenamesUnderANewKey(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\n")

	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo original"))

	pack := writePack(t, `[greet]
Name = "greet"
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	original := reloaded.AllCommandsByName.Get("greet")
	if original == nil || original.Exec != "echo original" {
		t.Errorf("original = %+v, want it preserved", original)
	}

	renamed := reloaded.AllCommandsByName.Get("greet_imported")
	if renamed == nil {
		t.Fatal("the renamed import is not in the configuration file")
	}

	if renamed.Name != "greet_imported" || renamed.Exec != "echo imported" {
		t.Errorf("renamed = %+v, want greet_imported carrying the pack's command", renamed)
	}
}

// A rename must not replace a command that already has the "_imported" name.
func TestImportRenameDoesNotReplaceAnExistingImportedCommand(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\n")

	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo original"))
	cfg.AllCommandsByName.Add("greet_imported", command.New("greet_imported", nil, "echo keep"))

	pack := writePack(t, `[greet]
Name = "greet"
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	if kept := reloaded.AllCommandsByName.Get("greet_imported"); kept == nil || kept.Exec != "echo keep" {
		t.Errorf("greet_imported = %+v, want it untouched", kept)
	}

	renamed := reloaded.AllCommandsByName.Get("greet_imported2")
	if renamed == nil || renamed.Exec != "echo imported" {
		t.Errorf("greet_imported2 = %+v, want the pack's command", renamed)
	}
}

func TestImportOverwriteReplacesTheCommand(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	original := command.New("greet", []string{"old"}, "echo original")
	cfg.AllCommandsByName.Add("greet", original)
	cfg.AllCommandsByAlias.Add("old", original)

	if err := sharedom.NewHandler(cfg).Import(writePack(t, greetPack), true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if cfg.AllCommandsByAlias.Exists("old") {
		t.Error("an alias of the replaced command survived in the running process")
	}

	reloaded := reload(t, cfgPath)

	replaced := reloaded.AllCommandsByName.Get("greet")
	if replaced == nil || replaced.Exec != "echo hi" {
		t.Errorf("greet = %+v, want the pack's command", replaced)
	}

	if reloaded.AllCommandsByAlias.Exists("old") {
		t.Error("the replaced command's alias survived the reload")
	}
}

// --overwrite may take a name that another command used as an alias, but the
// losing command must stop claiming it: otherwise the loaded map and the
// record disagree at the next run.
func TestImportOverwriteTakesANameThatIsAnAlias(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	owner := command.New("owner", []string{"greet"}, "echo owner")
	cfg.AllCommandsByName.Add("owner", owner)
	cfg.AllCommandsByAlias.Add("greet", owner)

	pack := writePack(t, `[greet]
Name = "greet"
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if got := cfg.AllCommandsByAlias.Get("greet"); got != nil {
		t.Errorf("alias greet still resolves to %+v after the name took it", got)
	}

	if got := cfg.AllCommandsByName.Get("owner"); got == nil || slices.Contains(got.AllAlias, "greet") {
		t.Errorf("owner = %+v, want it to drop the taken alias", got)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByName.Get("greet"); got == nil || got.Exec != "echo imported" {
		t.Errorf("greet = %+v, want the imported command", got)
	}

	if reloaded.AllCommandsByAlias.Exists("greet") {
		t.Error("the taken alias survived the reload")
	}
}

// A run that fails part-way must not leave a rewritten file behind, and must
// not leave half the file merged into the run.
func TestImportFailureWritesNothing(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\ny\n") // two answers, so the third prompt fails

	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo original"))
	cfg.AllCommandsByName.Add("alpha", command.New("alpha", nil, "echo alpha"))
	cfg.AllCommandsByName.Add("beta", command.New("beta", nil, "echo beta"))

	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("SaveCommands: %v", err)
	}

	before, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	pack := writePack(t, `[greet]
Name = "greet"
Exec = "echo imported"

[alpha]
Name = "alpha"
Exec = "echo imported alpha"

[beta]
Name = "beta"
Exec = "echo imported beta"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err == nil {
		t.Fatal("Import succeeded without an answer to every rename prompt")
	}

	after, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}

	if !bytes.Equal(before, after) {
		t.Errorf("a failed import rewrote the configuration file:\nbefore:\n%s\nafter:\n%s", before, after)
	}

	for _, key := range []string{
		"greet_imported", "greet_imported2",
		"alpha_imported", "alpha_imported2",
		"beta_imported", "beta_imported2",
	} {
		if cfg.AllCommandsByName.Exists(key) {
			t.Errorf("a failed import left %q in the run", key)
		}
	}
}

// Export keys the pack by the command's canonical name, so exporting by an
// alias cannot produce a pack that renames the command on import.
func TestExportKeysByCanonicalName(t *testing.T) {
	cfg, _ := newConfig(t, "")

	greet := command.New("greet", []string{"hi"}, "echo hi")
	cfg.AllCommandsByName.Add("greet", greet)
	cfg.AllCommandsByAlias.Add("hi", greet)

	outPath := filepath.Join(t.TempDir(), "pack.toml")
	if err := sharedom.NewHandler(cfg).Export("hi", outPath); err != nil {
		t.Fatalf("Export: %v", err)
	}

	imported, err := app.InitApp(outPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("read the exported pack: %v", err)
	}

	if imported.AllCommandsByName.Get("hi") != nil {
		t.Error("the pack is keyed by the alias")
	}

	if got := imported.AllCommandsByName.Get("greet"); got == nil || got.Exec != "echo hi" {
		t.Errorf("greet = %+v, want it keyed by its canonical name", got)
	}
}

// An alias collision is a name collision: the import must not let two commands
// claim the same key, or the alias resolves to whichever command the loader
// happens to walk last.
func TestImportRenamesAConflictingAlias(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\n")

	hello := command.New("hello", []string{"hi"}, "echo hello")
	cfg.AllCommandsByName.Add("hello", hello)
	cfg.AllCommandsByAlias.Add("hi", hello)

	pack := writePack(t, `[greet]
Name = "greet"
AllAlias = ["hi"]
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	// The running process must agree with what the next load will build.
	if got := cfg.AllCommandsByAlias.Get("hi"); got == nil || got.Name != "hello" {
		t.Errorf("alias hi = %+v, want it to stay on hello", got)
	}

	if got := cfg.AllCommandsByAlias.Get("hi_imported"); got == nil || got.Name != "greet" {
		t.Errorf("alias hi_imported = %+v, want the imported command", got)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByName.Get("greet"); got == nil || got.Exec != "echo imported" {
		t.Errorf("greet = %+v, want the imported command", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("hi"); got == nil || got.Name != "hello" {
		t.Errorf("reloaded alias hi = %+v, want hello", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("hi_imported"); got == nil || got.Name != "greet" {
		t.Errorf("reloaded alias hi_imported = %+v, want greet", got)
	}
}

// Declining the rename is a skip, not a partial import: the command and its
// conflicting alias stay out of both the run and the file.
func TestImportDecliningARenameSkipsTheCommand(t *testing.T) {
	cfg, cfgPath := newConfig(t, "n\n")

	hello := command.New("hello", []string{"hi"}, "echo hello")
	cfg.AllCommandsByName.Add("hello", hello)
	cfg.AllCommandsByAlias.Add("hi", hello)

	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("SaveCommands: %v", err)
	}

	pack := writePack(t, `[greet]
Name = "greet"
AllAlias = ["hi"]
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if cfg.AllCommandsByName.Exists("greet") {
		t.Error("the declined import reached the running process")
	}

	reloaded := reload(t, cfgPath)

	if reloaded.AllCommandsByName.Exists("greet") {
		t.Error("the declined import reached the configuration file")
	}

	if got := reloaded.AllCommandsByAlias.Get("hi"); got == nil || got.Name != "hello" {
		t.Errorf("alias hi = %+v, want hello", got)
	}
}

// A name that another command uses as an alias would shadow it on every
// name-first lookup, so the import renames the name.
func TestImportRenamesANameThatIsAnAlias(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\n")

	owner := command.New("owner", []string{"greet"}, "echo owner")
	cfg.AllCommandsByName.Add("owner", owner)
	cfg.AllCommandsByAlias.Add("greet", owner)

	pack := writePack(t, `[greet]
Name = "greet"
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByName.Get("greet"); got != nil {
		t.Errorf("the imported command took a name that was an alias: %+v", got)
	}

	if got := reloaded.AllCommandsByName.Get("greet_imported"); got == nil || got.Exec != "echo imported" {
		t.Errorf("greet_imported = %+v, want the imported command", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("greet"); got == nil || got.Name != "owner" {
		t.Errorf("alias greet = %+v, want it to stay on owner", got)
	}
}

// --overwrite does not prompt, and it does not take an alias from another
// command: the alias is dropped from the import instead.
func TestImportOverwriteDropsAConflictingAlias(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	hello := command.New("hello", []string{"hi"}, "echo hello")
	cfg.AllCommandsByName.Add("hello", hello)
	cfg.AllCommandsByAlias.Add("hi", hello)

	pack := writePack(t, `[greet]
Name = "greet"
AllAlias = ["hi"]
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByName.Get("greet"); got == nil || len(got.AllAlias) != 0 {
		t.Errorf("greet = %+v, want the command without the conflicting alias", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("hi"); got == nil || got.Name != "hello" {
		t.Errorf("alias hi = %+v, want it to stay on hello", got)
	}
}

// Re-importing the command that is being replaced must keep the aliases it
// already has: they belong to the same name, not to a stranger.
func TestImportOverwriteKeepsTheReplacedCommandsAliases(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	greet := command.New("greet", []string{"hi"}, "echo original")
	cfg.AllCommandsByName.Add("greet", greet)
	cfg.AllCommandsByAlias.Add("hi", greet)

	pack := writePack(t, `[greet]
Name = "greet"
AllAlias = ["hi"]
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByName.Get("greet"); got == nil || got.Exec != "echo imported" {
		t.Errorf("greet = %+v, want the imported exec", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("hi"); got == nil || got.Name != "greet" {
		t.Errorf("alias hi = %+v, want it on greet", got)
	}
}

// Renaming a conflicting alias must not land on the command's own name: that
// would create a self-alias the maps cannot distinguish.
func TestImportRenamedAliasDoesNotTakeTheCommandsName(t *testing.T) {
	cfg, _ := newConfig(t, "y\n")

	hello := command.New("hello", []string{"hi"}, "echo hello")
	cfg.AllCommandsByName.Add("hello", hello)
	cfg.AllCommandsByAlias.Add("hi", hello)

	// The conflicting alias "hi" wants "hi_imported", which is the new name.
	pack := writePack(t, `[hi_imported]
Name = "hi_imported"
AllAlias = ["hi"]
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if got := cfg.AllCommandsByAlias.Get("hi_imported"); got != nil {
		t.Errorf("alias hi_imported = %+v, want it to stay free for the name", got)
	}

	if got := cfg.AllCommandsByAlias.Get("hi_imported2"); got == nil || got.Name != "hi_imported" {
		t.Errorf("alias hi_imported2 = %+v, want the renamed alias", got)
	}
}

// Two entries in one pack share the namespace too: the second conflicting
// alias is renamed, not merged into the first command.
func TestImportResolvesCollisionsWithinThePack(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\n")

	pack := writePack(t, `[alpha]
Name = "alpha"
AllAlias = ["shared"]
Exec = "echo alpha"

[beta]
Name = "beta"
AllAlias = ["shared"]
Exec = "echo beta"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByAlias.Get("shared"); got == nil || got.Name != "alpha" {
		t.Errorf("alias shared = %+v, want alpha", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("shared_imported"); got == nil || got.Name != "beta" {
		t.Errorf("alias shared_imported = %+v, want beta", got)
	}
}

// --overwrite still has to produce distinct keys inside one pack; the later
// entry is renamed because there is nothing on disk to replace.
func TestImportOverwriteRenamesAPackInternalName(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	pack := writePack(t, `[a]
Name = "a"
AllAlias = ["b"]
Exec = "echo a"

[b]
Name = "b"
Exec = "echo b"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded := reload(t, cfgPath)

	if got := reloaded.AllCommandsByName.Get("b"); got != nil {
		t.Errorf("the later entry kept a name an earlier one uses as an alias: %+v", got)
	}

	if got := reloaded.AllCommandsByName.Get("b_imported"); got == nil || got.Exec != "echo b" {
		t.Errorf("b_imported = %+v, want the later entry", got)
	}

	if got := reloaded.AllCommandsByAlias.Get("b"); got == nil || got.Name != "a" {
		t.Errorf("alias b = %+v, want it on a", got)
	}
}

// Replacing a command must not delete an alias that a legacy duplicate has
// already handed to another command.
func TestImportOverwriteKeepsAnAliasAnotherCommandOwns(t *testing.T) {
	cfg, _ := newConfig(t, "")

	a := command.New("a", []string{"x"}, "echo a")
	cfg.AllCommandsByName.Add("a", a)

	b := command.New("b", []string{"x"}, "echo b")
	cfg.AllCommandsByName.Add("b", b)
	cfg.AllCommandsByAlias.Add("x", b)

	pack := writePack(t, `[a]
Name = "a"
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if got := cfg.AllCommandsByAlias.Get("x"); got == nil || got.Name != "b" {
		t.Errorf("alias x = %+v, want it to stay on b", got)
	}
}

// A hand-edited record without a Name can still own an alias; taking that key
// as a name must strip it from the record, not just from the alias map.
func TestImportOverwriteCleansAnOwnerWithoutAName(t *testing.T) {
	cfg, cfgPath := newConfig(t, "")

	owner := command.New("", []string{"greet"}, "echo owner")
	cfg.AllCommandsByName.Add("owner", owner)
	cfg.AllCommandsByAlias.Add("greet", owner)

	pack := writePack(t, `[greet]
Name = "greet"
Exec = "echo imported"
`)

	if err := sharedom.NewHandler(cfg).Import(pack, true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if got := cfg.AllCommandsByName.Get("owner"); got == nil || slices.Contains(got.AllAlias, "greet") {
		t.Errorf("owner = %+v, want it to drop the taken alias", got)
	}

	if reloaded := reload(t, cfgPath); reloaded.AllCommandsByAlias.Exists("greet") {
		t.Error("the taken alias came back after the reload")
	}
}
