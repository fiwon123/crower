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
	tag string
	err error
}

func (r *fakeReleases) LatestTag() (string, error) {
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

type fakeNotes struct{}

func (fakeNotes) UpgradeNote() string { return "note:upgrade" }

type harness struct {
	state    *fakeState
	output   *fakeOutput
	releases *fakeReleases
	updater  *fakeUpdater
}

func newHarness(tag string) *harness {
	return &harness{
		state:    &fakeState{},
		output:   &fakeOutput{},
		releases: &fakeReleases{tag: tag},
		updater:  &fakeUpdater{},
	}
}

func (h *harness) run(current string) error {
	return upgradeapp.New(h.state, h.output, h.releases, h.updater, fakeNotes{}).Run(current)
}

func TestUpgradeInstallsANewerRelease(t *testing.T) {
	h := newHarness("v0.9.0")

	if err := h.run("v0.8.0"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.updater.versions, []string{"v0.9.0"}; !slices.Equal(got, want) {
		t.Errorf("upgraded to %v, want %v", got, want)
	}

	want := []string{"Current Version: ", "crower upgraded from "}
	if len(h.output.infos) != 2 || h.output.infos[0] != want[0] || h.output.infos[1] != want[1] {
		t.Errorf("output = %v, want %v", h.output.infos, want)
	}

	if len(h.state.records) != 1 || h.state.records[0].op != operation.Upgrade || h.state.records[0].note != "note:upgrade" {
		t.Errorf("records = %+v, want the upgrade entry", h.state.records)
	}
}

func TestUpgradeAlreadyCurrent(t *testing.T) {
	h := newHarness("v0.8.0")

	if err := h.run("v0.8.0"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.updater.versions) != 0 {
		t.Errorf("upgraded %v, want nothing", h.updater.versions)
	}

	if len(h.state.records) != 0 {
		t.Errorf("records = %+v, want none", h.state.records)
	}

	if last := h.output.infos[len(h.output.infos)-1]; last != "already up-to-date" {
		t.Errorf("last diagnostic = %q, want already up-to-date", last)
	}
}

func TestUpgradeErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("latest tag", func(t *testing.T) {
		h := newHarness("")
		h.releases.err = boom

		if err := h.run("v0.8.0"); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("apply", func(t *testing.T) {
		h := newHarness("v0.9.0")
		h.updater.err = boom

		if err := h.run("v0.8.0"); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed upgrade recorded an entry")
		}
	})

	t.Run("record", func(t *testing.T) {
		h := newHarness("v0.9.0")
		h.state.recordErr = boom

		if err := h.run("v0.8.0"); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}
