package searchdom_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/searchdom"
	"github.com/fiwon123/crower/internal/helper"
)

// recordingRunner stands in for the shell: a search has to reach the runner it
// was given, not one it created for itself.
type recordingRunner struct {
	executed []string
	started  []string
	out      string
}

func (r *recordingRunner) Exec(ex string) (string, error) {
	r.executed = append(r.executed, ex)
	return r.out, nil
}

func (r *recordingRunner) Start(ex string) error {
	r.started = append(r.started, ex)
	return nil
}

func TestSearchUsesTheInjectedRunner(t *testing.T) {
	cfg, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("InitCrowerTests: %v", err)
	}

	t.Run("a filesystem search executes its command", func(t *testing.T) {
		runner := &recordingRunner{out: "/tmp/needle.txt\n"}
		handler := searchdom.NewHandlerWithRunner(cfg, runner)

		if _, err := handler.SearchFile(".", "needle.txt"); err != nil {
			t.Fatalf("SearchFile: %v", err)
		}

		if len(runner.executed) != 1 {
			t.Fatalf("executed %d commands, want 1: %v", len(runner.executed), runner.executed)
		}

		if !strings.Contains(runner.executed[0], "needle.txt") {
			t.Errorf("command = %q, want it to carry the search term", runner.executed[0])
		}
	})

	t.Run("a browser search starts its command instead", func(t *testing.T) {
		runner := &recordingRunner{}
		handler := searchdom.NewHandlerWithRunner(cfg, runner)

		if _, err := handler.SearchBrowser("golang"); err != nil {
			t.Fatalf("SearchBrowser: %v", err)
		}

		if len(runner.started) != 1 || len(runner.executed) != 0 {
			t.Fatalf("started %d, executed %d, want 1 started and 0 executed", len(runner.started), len(runner.executed))
		}

		if !strings.Contains(runner.started[0], "golang") {
			t.Errorf("command = %q, want it to carry the search term", runner.started[0])
		}
	})
}
