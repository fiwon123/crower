package checkrelease_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/checkrelease"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeReleases struct {
	tag string
	err error
}

func (r fakeReleases) LatestTag() (string, error) {
	return r.tag, r.err
}

type fakeOutput struct {
	lines []string
}

func (o *fakeOutput) Info(message string, data ...any) {
	line := message
	for _, value := range data {
		line += fmt.Sprint(value)
	}

	o.lines = append(o.lines, line)
}

func (o *fakeOutput) text() string {
	return strings.Join(o.lines, "\n")
}

type recorded struct {
	op          operation.Operation
	commandName string
	note        string
}

type fakeHistory struct {
	recorded *recorded
	err      error
}

func (h *fakeHistory) Record(op operation.Operation, commandName, note string) error {
	if h.err != nil {
		return h.err
	}

	h.recorded = &recorded{op: op, commandName: commandName, note: note}
	return nil
}

type fakeNotes struct{}

func (fakeNotes) Check() string { return "check-note" }

func newUseCase(releases checkrelease.Releases, output *fakeOutput, history *fakeHistory) *checkrelease.UseCase {
	return checkrelease.New(releases, output, history, fakeNotes{})
}

func TestRunReportsANewerRelease(t *testing.T) {
	output := &fakeOutput{}
	history := &fakeHistory{}

	if err := newUseCase(fakeReleases{tag: "v0.9.0"}, output, history).Run("v0.8.0"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(output.text(), "Current Version: ") || !strings.Contains(output.text(), "v0.8.0") {
		t.Errorf("the current version was not reported: %q", output.text())
	}

	if !strings.Contains(output.text(), "New Version Found: ") || !strings.Contains(output.text(), "v0.9.0") {
		t.Errorf("the new version was not reported: %q", output.text())
	}

	if history.recorded == nil {
		t.Fatal("a newer release was found but nothing was recorded")
	}

	if history.recorded.op != operation.Check {
		t.Errorf("recorded operation %q, want check", history.recorded.op)
	}

	if history.recorded.note != "check-note" {
		t.Errorf("recorded note %q, want the Notes port output", history.recorded.note)
	}
}

// Being up to date is not an event: it changes nothing, so it is not recorded.
func TestRunUpToDateRecordsNothing(t *testing.T) {
	output := &fakeOutput{}
	history := &fakeHistory{}

	if err := newUseCase(fakeReleases{tag: "v0.8.0"}, output, history).Run("v0.8.0"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(output.text(), "already up-to-date") {
		t.Errorf("the up-to-date case was not reported: %q", output.text())
	}

	if history.recorded != nil {
		t.Errorf("an up-to-date check recorded %+v", history.recorded)
	}
}

func TestRunPropagatesTheNetworkError(t *testing.T) {
	want := errors.New("network down")
	output := &fakeOutput{}
	history := &fakeHistory{}

	if err := newUseCase(fakeReleases{err: want}, output, history).Run("v0.8.0"); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}

	if history.recorded != nil {
		t.Fatal("a failed check recorded something")
	}
}

func TestRunPropagatesTheRecordError(t *testing.T) {
	want := errors.New("disk full")
	output := &fakeOutput{}
	history := &fakeHistory{err: want}

	if err := newUseCase(fakeReleases{tag: "v0.9.0"}, output, history).Run("v0.8.0"); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}
}
