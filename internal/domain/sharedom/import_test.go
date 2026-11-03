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

func writePack(t *testing.T, exec string) string {
	t.Helper()

	pack := filepath.Join(t.TempDir(), "pack.toml")
	contents := `[greet]
Name = "greet"
AllAlias = ["hi"]
Exec = "` + exec + `"
`

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
	pack := writePack(t, "echo hi")

	if err := sharedom.NewHandler(cfg).Import(pack, false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	// A second run reads the file, the way the next process would.
	reloaded, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("reload: %v", err)
	}

	imported := reloaded.AllCommandsByName.Get("greet")
	if imported == nil {
		t.Fatal("the imported command was not written to the configuration file")
	}

	if imported.Exec != "echo hi" {
		t.Errorf("exec = %q, want %q", imported.Exec, "echo hi")
	}

	if !slices.Contains(reloaded.OrderKeys, "greet") {
		t.Errorf("keys = %v, want the imported command", reloaded.OrderKeys)
	}
}

// Renaming has to insert the command under the new key, not overwrite the key
// it collided with while carrying a different Name.
func TestImportRenamesUnderANewKey(t *testing.T) {
	cfg, cfgPath := newConfig(t, "y\n")

	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo original"))

	if err := sharedom.NewHandler(cfg).Import(writePack(t, "echo imported"), false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	reloaded, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("reload: %v", err)
	}

	original := reloaded.AllCommandsByName.Get("greet")
	if original == nil || original.Exec != "echo original" {
		t.Errorf("original = %+v, want it preserved", original)
	}

	renamed := reloaded.AllCommandsByName.Get("greet_imported")
	if renamed == nil {
		t.Fatal("the renamed import is not in the configuration file")
	}

	if renamed.Name != "greet_imported" {
		t.Errorf("name = %q, want greet_imported", renamed.Name)
	}

	if renamed.Exec != "echo imported" {
		t.Errorf("exec = %q, want the imported command", renamed.Exec)
	}
}

// A run that fails part-way must not leave a rewritten file behind.
func TestImportFailureWritesNothing(t *testing.T) {
	cfg, cfgPath := newConfig(t, "") // a closed input stream fails the prompt

	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo original"))
	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("SaveCommands: %v", err)
	}

	before, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	if err := sharedom.NewHandler(cfg).Import(writePack(t, "echo imported"), false); err == nil {
		t.Fatal("Import succeeded without an answer to the rename prompt")
	}

	after, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}

	if !bytes.Equal(before, after) {
		t.Errorf("a failed import rewrote the configuration file:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
