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

	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo original"))

	if err := sharedom.NewHandler(cfg).Import(writePack(t, greetPack), true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	replaced := reload(t, cfgPath).AllCommandsByName.Get("greet")
	if replaced == nil || replaced.Exec != "echo hi" {
		t.Errorf("greet = %+v, want the pack's command", replaced)
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
