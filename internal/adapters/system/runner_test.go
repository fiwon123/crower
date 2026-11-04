package system

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	crerrors "github.com/fiwon123/crower/internal/errors"
)

// The nesting level comes from the environment, because the recursive call is a
// new crower process: an in-process counter never sees it.
func TestNestedDepthReadsTheEnvironment(t *testing.T) {
	t.Setenv(nestedDepthEnv, "3")

	if got := nestedDepth(); got != 3 {
		t.Errorf("nestedDepth() = %d, want 3", got)
	}
}

func TestNestedDepthIsZeroWithoutTheVariable(t *testing.T) {
	t.Setenv(nestedDepthEnv, "")

	if got := nestedDepth(); got != 0 {
		t.Errorf("nestedDepth() = %d, want 0", got)
	}
}

// A command whose steps call itself would spawn a process per step forever. The
// limit is the only thing standing between that command and the machine running
// out of processes.
func TestRunNestedRefusesToGoTooDeep(t *testing.T) {
	t.Setenv(nestedDepthEnv, strconv.Itoa(maxExecDepth))

	_, err := NewRunner().RunNested("crower list", Target{})
	if err == nil {
		t.Fatal("RunNested() at the maximum depth returned no error")
	}

	if !strings.Contains(err.Error(), "maximum depth") {
		t.Errorf("error = %v, want it to name the limit", err)
	}
}

// An empty Cwd is the process's own directory, not a folder named "".
func TestResolveDirEmpty(t *testing.T) {
	got, err := resolveDir("")
	if err != nil || got != "" {
		t.Fatalf("resolveDir(\"\") = %q, %v; want empty and no error", got, err)
	}
}

func TestResolveDirAcceptsADirectory(t *testing.T) {
	dir := t.TempDir()

	got, err := resolveDir(dir)
	if err != nil {
		t.Fatalf("resolveDir(%q): %v", dir, err)
	}

	if got != dir {
		t.Errorf("resolveDir(%q) = %q, want it unchanged", dir, got)
	}
}

// "~" is the user's home, and "~/x" is a folder below it.
func TestResolveDirExpandsHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if err := os.Mkdir(filepath.Join(home, "projects"), 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	tests := []struct {
		name string
		dir  string
		want string
	}{
		{"bare tilde", "~", home},
		{"tilde slash", "~/projects", filepath.Join(home, "projects")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveDir(tt.dir)
			if err != nil {
				t.Fatalf("resolveDir(%q): %v", tt.dir, err)
			}

			if got != tt.want {
				t.Errorf("resolveDir(%q) = %q, want %q", tt.dir, got, tt.want)
			}
		})
	}
}

// A missing folder is refused before a process starts, with a sentinel a caller
// can test rather than a message it has to parse.
func TestResolveDirRejectsAMissingFolder(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	_, err := resolveDir(missing)
	if !errors.Is(err, crerrors.ErrInvalidWorkDir) {
		t.Fatalf("error = %v, want ErrInvalidWorkDir", err)
	}
}

// A file is not a working directory.
func TestResolveDirRejectsAFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err := resolveDir(path)
	if !errors.Is(err, crerrors.ErrInvalidWorkDir) {
		t.Fatalf("error = %v, want ErrInvalidWorkDir", err)
	}
}

// Only the current user's home is expanded; "~other/x" is a relative path, not
// a lookup of another account.
func TestExpandHomeLeavesAnotherUserAlone(t *testing.T) {
	got, err := expandHome("~other/x")
	if err != nil {
		t.Fatalf("expandHome: %v", err)
	}

	if got != "~other/x" {
		t.Errorf("expandHome(\"~other/x\") = %q, want it unchanged", got)
	}
}

// A literal value is accepted unless its key looks like a credential.
func TestBuildEnvOverridesAndSorts(t *testing.T) {
	t.Setenv("CROWER_TEST_ORDER", "process")

	env, err := buildEnv(map[string]string{"CROWER_TEST_ORDER": "override", "CROWER_TEST_A": "1"})
	if err != nil {
		t.Fatalf("buildEnv: %v", err)
	}

	count := 0

	for _, entry := range env {
		if strings.HasPrefix(entry, "CROWER_TEST_ORDER=") {
			count++

			if entry != "CROWER_TEST_ORDER=override" {
				t.Errorf("entry = %q, want the override", entry)
			}
		}
	}

	if count != 1 {
		t.Errorf("CROWER_TEST_ORDER appears %d times, want once", count)
	}

	if !slices.IsSorted(env) {
		t.Errorf("environment is not sorted: %v", env)
	}
}

// A malformed or non-positive timeout is refused before anything runs.
func TestTimeoutContextRejectsBadDurations(t *testing.T) {
	for _, raw := range []string{"soon", "0s", "-1s", "5"} {
		if _, _, err := timeoutContext(raw); !errors.Is(err, crerrors.ErrInvalidTimeout) {
			t.Errorf("timeoutContext(%q) error = %v, want ErrInvalidTimeout", raw, err)
		}
	}
}
