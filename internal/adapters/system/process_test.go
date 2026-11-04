package system

import (
	"errors"
	"testing"

	"github.com/fiwon123/crower/internal/shell"
)

// stubProcesses builds a Processes with every lookup answered by the test, so
// no real process table is needed.
func stubProcesses() *Processes {
	return &Processes{
		pathByID:   func(int32) (string, error) { return "", nil },
		nameByID:   func(int32) (string, error) { return "", nil },
		pathByName: func(string) (string, error) { return "", nil },
		flatpakID:  func(string) (string, error) { return "", nil },
	}
}

func TestExecLineByNameQuotesThePath(t *testing.T) {
	p := stubProcesses()
	p.pathByName = func(name string) (string, error) {
		if name != "firefox" {
			t.Errorf("looked up %q, want firefox", name)
		}

		return "/usr/bin/firefox", nil
	}

	got, err := p.ExecLine("firefox")
	if err != nil {
		t.Fatalf("ExecLine: %v", err)
	}

	if want := shell.Quote("/usr/bin/firefox"); got != want {
		t.Errorf("ExecLine = %q, want %q", got, want)
	}
}

func TestExecLineByPIDQuotesThePath(t *testing.T) {
	p := stubProcesses()
	p.pathByID = func(pid int32) (string, error) {
		if pid != 4242 {
			t.Errorf("looked up pid %d, want 4242", pid)
		}

		return "/opt/My App/bin/app", nil
	}

	got, err := p.ExecLine("4242")
	if err != nil {
		t.Fatalf("ExecLine: %v", err)
	}

	if want := shell.Quote("/opt/My App/bin/app"); got != want {
		t.Errorf("ExecLine = %q, want %q", got, want)
	}
}

// A flatpak executable lives inside a sandbox path; the line a user would type
// is `flatpak run <app id>`, so the app id is looked up by the process name.
func TestExecLineByPIDResolvesAFlatpakApp(t *testing.T) {
	p := stubProcesses()
	p.pathByID = func(int32) (string, error) { return "/app/bin/firefox", nil }
	p.nameByID = func(int32) (string, error) { return "firefox", nil }
	p.flatpakID = func(name string) (string, error) {
		if name != "firefox" {
			t.Errorf("flatpak lookup for %q, want firefox", name)
		}

		return "org.mozilla.firefox", nil
	}

	got, err := p.ExecLine("7")
	if err != nil {
		t.Fatalf("ExecLine: %v", err)
	}

	if want := "flatpak run org.mozilla.firefox"; got != want {
		t.Errorf("ExecLine = %q, want %q", got, want)
	}
}

// A name lookup already knows the process name, so it must not ask for it
// again -- the pid it would need came from parsing a name, which failed.
func TestExecLineByNameResolvesAFlatpakAppWithoutAPIDLookup(t *testing.T) {
	p := stubProcesses()
	p.pathByName = func(string) (string, error) { return "/app/bin/firefox", nil }
	p.nameByID = func(int32) (string, error) {
		t.Error("nameByID was called for a name lookup")
		return "", nil
	}
	p.flatpakID = func(string) (string, error) { return "org.mozilla.firefox", nil }

	if _, err := p.ExecLine("firefox"); err != nil {
		t.Fatalf("ExecLine: %v", err)
	}
}

func TestExecLineWithAnUnknownProcessFails(t *testing.T) {
	p := stubProcesses()

	if _, err := p.ExecLine("ghost"); err == nil || err.Error() != "couldn't find the process either by pid or by name" {
		t.Fatalf("error = %v, want the not-found message", err)
	}
}

func TestExecLinePassesLookupErrorsUp(t *testing.T) {
	boom := errors.New("boom")

	t.Run("by name", func(t *testing.T) {
		p := stubProcesses()
		p.pathByName = func(string) (string, error) { return "", boom }

		if _, err := p.ExecLine("ghost"); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("flatpak", func(t *testing.T) {
		p := stubProcesses()
		p.pathByID = func(int32) (string, error) { return "/app/bin/firefox", nil }
		p.nameByID = func(int32) (string, error) { return "firefox", nil }
		p.flatpakID = func(string) (string, error) { return "", boom }

		if _, err := p.ExecLine("7"); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}
