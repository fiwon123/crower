package reverthistory_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/reverthistory"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeHistory struct {
	dropped  string
	kept     string
	pending  int
	version  int
	err      error
	recorded *struct {
		op          operation.Operation
		commandName string
		note        string
	}
}

func (h *fakeHistory) Pending(steps int) (string, string, error) {
	h.pending = steps

	if h.err != nil {
		return "", "", h.err
	}

	return h.dropped, h.kept, nil
}

func (h *fakeHistory) RevertTo(int) (int, error) {
	if h.err != nil {
		return 0, h.err
	}

	return h.version, nil
}

func (h *fakeHistory) Record(op operation.Operation, commandName, note string) error {
	h.recorded = &struct {
		op          operation.Operation
		commandName string
		note        string
	}{op, commandName, note}

	return nil
}

type fakeOutput struct{ lines []string }

func (o *fakeOutput) Info(message string, data ...any) {
	line := message
	for _, value := range data {
		line += fmt.Sprint(value)
	}

	o.lines = append(o.lines, line)
}

func (o *fakeOutput) text() string { return strings.Join(o.lines, "\n") }

type fakePrompt struct {
	answer bool
	err    error
	asked  []string
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)
	return p.answer, p.err
}

type fakeNotes struct{}

func (fakeNotes) Revert(args []string) string {
	return "revert-note:" + strings.Join(args, ",")
}

func newUseCase(history reverthistory.History, output *fakeOutput, prompt *fakePrompt) *reverthistory.UseCase {
	return reverthistory.New(history, output, prompt, fakeNotes{})
}

func TestRunRevertsAndRecords(t *testing.T) {
	history := &fakeHistory{dropped: "DROPPED", kept: "KEPT", version: 3}
	output := &fakeOutput{}
	prompt := &fakePrompt{answer: true}

	if err := newUseCase(history, output, prompt).Run([]string{"2"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if history.pending != 2 {
		t.Errorf("pending asked for %d steps, want 2", history.pending)
	}

	for _, want := range []string{"Deleted History", "DROPPED", "New History", "KEPT", "History will revert 2"} {
		if !strings.Contains(output.text(), want) {
			t.Errorf("diagnostics do not contain %q: %q", want, output.text())
		}
	}

	if !strings.Contains(output.text(), "reverted to history version 3") {
		t.Errorf("the new version was not reported: %q", output.text())
	}

	if history.recorded == nil {
		t.Fatal("the revert was not recorded")
	}

	if history.recorded.op != operation.Revert {
		t.Errorf("recorded operation %q, want revert", history.recorded.op)
	}

	if history.recorded.note != "revert-note:2" {
		t.Errorf("recorded note %q, want the Notes port output", history.recorded.note)
	}
}

// Declining is not a failure, and nothing is recorded.
func TestRunDeclines(t *testing.T) {
	history := &fakeHistory{dropped: "D", kept: "K"}
	output := &fakeOutput{}
	prompt := &fakePrompt{answer: false}

	if err := newUseCase(history, output, prompt).Run([]string{"1"}); err != nil {
		t.Fatalf("a declined revert returned an error: %v", err)
	}

	if history.recorded != nil {
		t.Fatal("a declined revert recorded something")
	}

	if !strings.Contains(output.text(), "Cancelling revert...") {
		t.Errorf("the decline was not reported: %q", output.text())
	}
}

func TestRunValidatesTheStepCount(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no argument", nil, "missing arguments"},
		{"not a number", []string{"two"}, "must be a number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := &fakeHistory{}
			output := &fakeOutput{}
			prompt := &fakePrompt{answer: true}

			err := newUseCase(history, output, prompt).Run(tt.args)
			if err == nil {
				t.Fatal("Run returned no error")
			}

			if !strings.Contains(strings.ToLower(err.Error()), tt.want) {
				t.Errorf("error = %v, want it to mention %q", err, tt.want)
			}

			if history.recorded != nil {
				t.Fatal("an invalid step count recorded something")
			}
		})
	}
}

// An out-of-range step count is refused before the user is asked anything.
func TestRunRefusesAnOutOfRangeStepCount(t *testing.T) {
	want := errors.New("can't revert: steps is greater than quantity of history registry")
	history := &fakeHistory{err: want}
	output := &fakeOutput{}
	prompt := &fakePrompt{answer: true}

	if err := newUseCase(history, output, prompt).Run([]string{"99"}); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}

	if len(prompt.asked) != 0 {
		t.Errorf("asked %v before validating the range", prompt.asked)
	}

	if strings.Contains(output.text(), "Deleted History") {
		t.Errorf("the lists were printed for an invalid revert: %q", output.text())
	}
}

func TestRunPropagatesThePromptError(t *testing.T) {
	want := errors.New("input is required")
	history := &fakeHistory{}
	output := &fakeOutput{}
	prompt := &fakePrompt{err: want}

	if err := newUseCase(history, output, prompt).Run([]string{"1"}); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}
}
