package upgradeapp_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/fiwon123/crower/internal/application/upgradeapp"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeState struct {
	records   []record
	recordErr error
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

type fakeOutput struct{ infos []string }

func (o *fakeOutput) Info(message string, data ...any) {
	o.infos = append(o.infos, message)
}

type fakeReleases struct {
	tag   string
	err   error
	calls int
}

func (r *fakeReleases) LatestTag() (string, error) {
	r.calls++
	return r.tag, r.err
}

type fakeUpdater struct {
	versions []string
	err      error
}

func (u *fakeUpdater) Upgrade(version string) error {
	u.versions = append(u.versions, version)
	return u.err
}

type fakeVersions struct {
	installed []string
	recordErr error
}

func (v *fakeVersions) Previous(current string) (string, bool) {
	for _, version := range v.installed {
		if version != current {
			return version, true
		}
	}

	return "", false
}

func (v *fakeVersions) RecordInstalled(version string) error {
	if v.recordErr != nil {
		return v.recordErr
	}

	v.installed = append([]string{version}, v.installed...)
	return nil
}

type fakePrompt struct {
	answer     bool
	err        error
	asked      []string
	confirmErr error
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)

	if p.confirmErr != nil {
		return false, p.confirmErr
	}

	return p.answer, p.err
}

type fakeNotes struct{}

func (fakeNotes) UpgradeNote() string { return "note:upgrade" }

type harness struct {
	state    *fakeState
	output   *fakeOutput
	releases *fakeReleases
	updater  *fakeUpdater
	versions *fakeVersions
	prompt   *fakePrompt
}

func newHarness(tag string) *harness {
	return &harness{
		state:    &fakeState{},
		output:   &fakeOutput{},
		releases: &fakeReleases{tag: tag},
		updater:  &fakeUpdater{},
		versions: &fakeVersions{},
		prompt:   &fakePrompt{answer: true},
	}
}

func (h *harness) run(request upgradeapp.Request) error {
	return upgradeapp.New(
		h.state, h.output, h.releases, h.updater, h.versions, h.prompt, fakeNotes{},
	).Run(request)
}

func TestUpdateInstallsTheNewestOnAYes(t *testing.T) {
	h := newHarness("v0.9.0")

	if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.updater.versions, []string{"v0.9.0"}; !slices.Equal(got, want) {
		t.Errorf("installed %v, want %v", got, want)
	}

	if got, want := h.prompt.asked, []string{"Install v0.9.0?"}; !slices.Equal(got, want) {
		t.Errorf("questions = %v, want %v", got, want)
	}

	if got, want := h.versions.installed, []string{"v0.9.0"}; !slices.Equal(got, want) {
		t.Errorf("recorded versions = %v, want %v", got, want)
	}

	if len(h.state.records) != 1 || h.state.records[0].op != operation.Upgrade {
		t.Errorf("records = %+v, want the upgrade entry", h.state.records)
	}

	if !slices.Contains(h.output.infos, "crower installed ") {
		t.Errorf("output = %v, want the install line", h.output.infos)
	}
}

func TestUpdateAlreadyCurrentAsksNothing(t *testing.T) {
	h := newHarness("v0.8.0")

	if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.updater.versions) != 0 || len(h.prompt.asked) != 0 || len(h.state.records) != 0 {
		t.Errorf("an up-to-date run did work: %v %v %v", h.updater.versions, h.prompt.asked, h.state.records)
	}

	if last := h.output.infos[len(h.output.infos)-1]; last != "already up-to-date" {
		t.Errorf("last diagnostic = %q, want already up-to-date", last)
	}
}

func TestUpdateDeclined(t *testing.T) {
	h := newHarness("v0.9.0")
	h.prompt.answer = false

	if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.updater.versions) != 0 || len(h.state.records) != 0 {
		t.Error("a declined update installed or recorded something")
	}

	if !slices.Contains(h.output.infos, "Cancelling update...") {
		t.Errorf("output = %v, want the cancel line", h.output.infos)
	}
}

// A named version installs even when it is older than the running one, and the
// release check is not consulted.
func TestUpdateInstallsANamedVersion(t *testing.T) {
	h := newHarness("")
	h.releases.err = errors.New("the release check must not run")

	if err := h.run(upgradeapp.Request{Current: "v0.8.0", Version: "v0.7.0"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.updater.versions, []string{"v0.7.0"}; !slices.Equal(got, want) {
		t.Errorf("installed %v, want %v", got, want)
	}

	if h.releases.calls != 0 {
		t.Errorf("LatestTag called %d time(s), want none", h.releases.calls)
	}

	if got, want := h.prompt.asked, []string{"Install v0.7.0?"}; !slices.Equal(got, want) {
		t.Errorf("questions = %v, want %v", got, want)
	}
}

func TestUpdateRollback(t *testing.T) {
	h := newHarness("")
	h.releases.err = errors.New("the release check must not run")
	h.versions.installed = []string{"v0.8.0", "v0.7.0"}

	if err := h.run(upgradeapp.Request{Current: "v0.8.0", Rollback: true}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.updater.versions, []string{"v0.7.0"}; !slices.Equal(got, want) {
		t.Errorf("installed %v, want %v", got, want)
	}

	if h.releases.calls != 0 {
		t.Errorf("LatestTag called %d time(s), want none", h.releases.calls)
	}
}

func TestUpdateRollbackWithoutHistory(t *testing.T) {
	h := newHarness("")

	err := h.run(upgradeapp.Request{Current: "v0.8.0", Rollback: true})
	if err == nil {
		t.Fatal("a rollback with no recorded version returned no error")
	}

	if len(h.updater.versions) != 0 {
		t.Error("a failed rollback installed something")
	}
}

func TestUpdateAssumeYesSkipsTheQuestion(t *testing.T) {
	h := newHarness("v0.9.0")

	if err := h.run(upgradeapp.Request{Current: "v0.8.0", AssumeYes: true}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.prompt.asked) != 0 {
		t.Errorf("asked %v, want no question", h.prompt.asked)
	}

	if got, want := h.updater.versions, []string{"v0.9.0"}; !slices.Equal(got, want) {
		t.Errorf("installed %v, want %v", got, want)
	}
}

func TestUpdateErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("latest tag", func(t *testing.T) {
		h := newHarness("")
		h.releases.err = boom

		if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("install", func(t *testing.T) {
		h := newHarness("v0.9.0")
		h.updater.err = boom

		if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.versions.installed) != 0 || len(h.state.records) != 0 {
			t.Error("a failed install recorded something")
		}
	})

	t.Run("recorded version", func(t *testing.T) {
		h := newHarness("v0.9.0")
		h.versions.recordErr = boom

		if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed version record still wrote a history entry")
		}
	})

	t.Run("prompt", func(t *testing.T) {
		h := newHarness("v0.9.0")
		h.prompt.confirmErr = boom

		if err := h.run(upgradeapp.Request{Current: "v0.8.0"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}
