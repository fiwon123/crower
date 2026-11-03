package sharedom_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/sharedom"
	crowererrors "github.com/fiwon123/crower/internal/errors"
)

// newConfigWithInput is newConfig with a reader the test controls, so a helper
// can observe the world while a prompt is pending.
func newConfigWithInput(t *testing.T, in io.Reader) (*app.Config, string) {
	t.Helper()

	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	cfg, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: in, Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	return cfg, cfgPath
}

func newPackServer(t *testing.T, contents string) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/packs/good.toml", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, contents)
	})

	return httptest.NewTLSServer(mux)
}

// tempSpy records the temp directory every time the prompt reads an answer, so
// a test can see what existed while the confirmation was pending.
type tempSpy struct {
	tmpDir string
	reader io.Reader
	seen   []string
}

func (s *tempSpy) Read(p []byte) (int, error) {
	entries, err := os.ReadDir(s.tmpDir)
	if err == nil {
		for _, entry := range entries {
			s.seen = append(s.seen, entry.Name())
		}
	}

	return s.reader.Read(p)
}

func assertNoPackTempFiles(t *testing.T, dir string) {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read temp dir: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "crower-pack-") || entry.Name() == "good.toml" {
			t.Errorf("pack temp file left behind: %s", entry.Name())
		}
	}
}

// A traversing name is rejected before any request: the name would otherwise
// become part of a path.
func TestPullRejectsAPathName(t *testing.T) {
	cfg, _ := newConfig(t, "")

	reached := false
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		http.NotFound(w, r)
	}))
	defer server.Close()

	handler := sharedom.NewHandler(cfg).WithRegistry(server.URL, server.Client())

	for _, name := range []string{"../escape", "a/b", `a\b`, "..", "", "good?x", "café"} {
		err := handler.Pull(name)
		if !errors.Is(err, crowererrors.ErrInvalidPackName) {
			t.Errorf("Pull(%q) error = %v, want ErrInvalidPackName", name, err)
		}
	}

	if reached {
		t.Error("an invalid pack name reached the network")
	}
}

func TestPullRejectsANonHTTPSRegistry(t *testing.T) {
	cfg, _ := newConfig(t, "")

	err := sharedom.NewHandler(cfg).WithRegistry("http://example.com", nil).Pull("good")
	if err == nil || !strings.Contains(err.Error(), "only HTTPS") {
		t.Errorf("error = %v, want the HTTPS refusal", err)
	}
}

// Nothing is written before the confirmation. The spy reads the temp directory
// while the prompt is pending: the old code had already written the pack there.
func TestPullWritesNothingBeforeTheConfirmation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.TempDir does not follow TMPDIR on Windows")
	}

	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)

	spy := &tempSpy{tmpDir: tmp, reader: strings.NewReader("n\n")}
	cfg, cfgPath := newConfigWithInput(t, spy)

	server := newPackServer(t, greetPack)
	defer server.Close()

	handler := sharedom.NewHandler(cfg).WithRegistry(server.URL, server.Client())

	if err := handler.Pull("good"); err != nil {
		t.Fatalf("Pull: %v", err)
	}

	if len(spy.seen) != 0 {
		t.Errorf("files existed while the confirmation was pending: %v", spy.seen)
	}

	assertNoPackTempFiles(t, tmp)

	if cfg.AllCommandsByName.Exists("greet") {
		t.Error("a declined pull imported the command")
	}

	if reload(t, cfgPath).AllCommandsByName.Exists("greet") {
		t.Error("a declined pull reached the configuration file")
	}
}

// Accepted: the pack is imported from a program-named temp file, which is
// removed afterwards.
func TestPullImportsAfterConfirmation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.TempDir does not follow TMPDIR on Windows")
	}

	cfg, cfgPath := newConfig(t, "y\n")

	server := newPackServer(t, greetPack)
	defer server.Close()

	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)

	handler := sharedom.NewHandler(cfg).WithRegistry(server.URL, server.Client())

	if err := handler.Pull("good"); err != nil {
		t.Fatalf("Pull: %v", err)
	}

	if got := reload(t, cfgPath).AllCommandsByName.Get("greet"); got == nil || got.Exec != "echo hi" {
		t.Errorf("greet = %+v, want the pulled command", got)
	}

	assertNoPackTempFiles(t, tmp)
}

// An empty pack is a failure, not a prompt that imports nothing.
func TestPullRejectsAnEmptyPack(t *testing.T) {
	cfg, _ := newConfig(t, "")

	server := newPackServer(t, "")
	defer server.Close()

	err := sharedom.NewHandler(cfg).WithRegistry(server.URL, server.Client()).Pull("good")
	if err == nil || !strings.Contains(err.Error(), "is empty") {
		t.Errorf("error = %v, want an empty-pack error", err)
	}
}

// A pack carries the user's commands, sometimes secrets. The log file must show
// the entries, never the body.
func TestPullDoesNotLogThePackBody(t *testing.T) {
	cfg, cfgPath := newConfig(t, "n\n")

	server := newPackServer(t, greetPack)
	defer server.Close()

	handler := sharedom.NewHandler(cfg).WithRegistry(server.URL, server.Client())

	if err := handler.Pull("good"); err != nil {
		t.Fatalf("Pull: %v", err)
	}

	logged, err := os.ReadFile(filepath.Join(filepath.Dir(cfgPath), "crower.log"))
	if err != nil {
		t.Fatalf("read the log: %v", err)
	}

	if strings.Contains(string(logged), "echo hi") {
		t.Errorf("the pack body reached the log file:\n%s", logged)
	}

	if !strings.Contains(string(logged), "greet") {
		t.Errorf("the log does not show the pack's entry names:\n%s", logged)
	}
}
