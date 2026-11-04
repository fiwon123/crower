package sharecommands_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/sharecommands"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

type fakeState struct {
	commands catalog.Catalog
	saved    []catalog.Catalog
	saveErr  error
}

func (s *fakeState) Commands() catalog.Catalog {
	return s.commands
}

func (s *fakeState) Save(commands catalog.Catalog) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.saved = append(s.saved, commands)
	return nil
}

type fakeOutput struct {
	infos   []string
	results []string
}

func (o *fakeOutput) Info(message string, data ...any) {
	o.infos = append(o.infos, fmt.Sprint(append([]any{message}, data...)...))
}

func (o *fakeOutput) Result(text string) {
	o.results = append(o.results, text)
}

type fakePrompt struct {
	answers []bool
	err     error
	asked   []string
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)

	if p.err != nil {
		return false, p.err
	}

	if len(p.answers) == 0 {
		return false, nil
	}

	answer := p.answers[0]
	p.answers = p.answers[1:]

	return answer, nil
}

type fakePacks struct {
	written  map[string][]catalog.Command
	read     []catalog.Command
	readErr  error
	writeErr error
}

func (p *fakePacks) Write(path string, commands []catalog.Command) error {
	if p.writeErr != nil {
		return p.writeErr
	}

	if p.written == nil {
		p.written = map[string][]catalog.Command{}
	}

	p.written[path] = commands
	return nil
}

func (p *fakePacks) Read(string) ([]catalog.Command, error) {
	return p.read, p.readErr
}

type fakeRegistry struct {
	packs    []sharecommands.Pack
	pack     []catalog.Command
	url      string
	indexErr error
	packErr  error
}

func (r *fakeRegistry) Index() ([]sharecommands.Pack, error) {
	return r.packs, r.indexErr
}

func (r *fakeRegistry) Pack(string) ([]catalog.Command, string, error) {
	return r.pack, r.url, r.packErr
}

type harness struct {
	state  *fakeState
	output *fakeOutput
	prompt *fakePrompt
	packs  *fakePacks
	reg    *fakeRegistry
}

func newHarness(commands ...catalog.Command) *harness {
	return &harness{
		state:  &fakeState{commands: catalog.New(nil, commands)},
		output: &fakeOutput{},
		prompt: &fakePrompt{},
		packs:  &fakePacks{},
		reg:    &fakeRegistry{},
	}
}

func (h *harness) useCase() *sharecommands.UseCase {
	return sharecommands.New(h.state, h.output, h.prompt, h.packs, h.reg)
}

func TestExportAllAndOne(t *testing.T) {
	h := newHarness(
		catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"},
		catalog.Command{Name: "build", Exec: "go build"},
	)
	useCase := h.useCase()

	if err := useCase.Export("", ""); err != nil {
		t.Fatalf("Export all: %v", err)
	}

	if _, ok := h.packs.written["crower-export.toml"]; !ok {
		t.Errorf("wrote %v, want the default path", h.packs.written)
	}

	if err := useCase.Export("hi", "one.toml"); err != nil {
		t.Fatalf("Export one: %v", err)
	}

	written := h.packs.written["one.toml"]
	if len(written) != 1 || written[0].Name != "greet" {
		t.Errorf("exported %+v, want greet by its alias", written)
	}

	if err := useCase.Export("ghost", ""); err == nil || err.Error() != "command not found" {
		t.Fatalf("error = %v, want command not found", err)
	}
}

func TestImportPlainPack(t *testing.T) {
	h := newHarness()
	h.packs.read = []catalog.Command{{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}}

	if err := h.useCase().Import("pack.toml", false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if len(h.state.saved) != 1 {
		t.Fatalf("saved %d catalogs, want 1", len(h.state.saved))
	}

	got, found := h.state.saved[0].Get("greet")
	if !found || got.Exec != "echo hi" {
		t.Errorf("saved = %+v, %v; want the imported command", got, found)
	}

	if owner, _ := h.state.saved[0].AliasOwner("hi"); owner != "greet" {
		t.Errorf("alias hi = %q, want greet", owner)
	}
}

func TestImportRenamesAConflictingNameWhenConfirmed(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Exec: "echo local"})
	h.packs.read = []catalog.Command{{Name: "greet", Exec: "echo imported"}}
	h.prompt.answers = []bool{true}

	if err := h.useCase().Import("pack.toml", false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if _, found := h.state.saved[0].Get("greet_imported"); !found {
		t.Errorf("names = %v, want greet_imported", h.state.saved[0].Names())
	}

	if !strings.Contains(h.prompt.asked[0], "Command 'greet' conflicts") {
		t.Errorf("question = %q, want the conflict prompt", h.prompt.asked[0])
	}

	// The local command is untouched.
	local, _ := h.state.saved[0].Get("greet")
	if local.Exec != "echo local" {
		t.Errorf("local exec = %q, want it kept", local.Exec)
	}
}

func TestImportSkipsAConflictWhenDeclined(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Exec: "echo local"})
	h.packs.read = []catalog.Command{{Name: "greet", Exec: "echo imported"}}
	h.prompt.answers = []bool{false}

	if err := h.useCase().Import("pack.toml", false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if len(h.state.saved) != 0 {
		t.Errorf("saved %v, want nothing", h.state.saved)
	}
}

func TestImportOverwriteReplaces(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Aliases: []string{"old"}, Exec: "echo local"})
	h.packs.read = []catalog.Command{{Name: "greet", Aliases: []string{"new"}, Exec: "echo imported"}}

	if err := h.useCase().Import("pack.toml", true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	got, _ := h.state.saved[0].Get("greet")
	if got.Exec != "echo imported" {
		t.Errorf("exec = %q, want the imported line", got.Exec)
	}

	if owner, found := h.state.saved[0].AliasOwner("old"); found {
		t.Errorf("alias old still resolves to %q, want it dropped", owner)
	}

	if owner, _ := h.state.saved[0].AliasOwner("new"); owner != "greet" {
		t.Errorf("alias new = %q, want greet", owner)
	}

	if len(h.prompt.asked) != 0 {
		t.Errorf("overwrite asked %v, want no questions", h.prompt.asked)
	}
}

// Re-importing the same command keeps the aliases it already owns under
// --overwrite: they are not a conflict with itself.
func TestImportOverwriteKeepsItsOwnAliases(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo local"})
	h.packs.read = []catalog.Command{{Name: "greet", Aliases: []string{"hi"}, Exec: "echo imported"}}

	if err := h.useCase().Import("pack.toml", true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	got, _ := h.state.saved[0].Get("greet")
	if !slices.Equal(got.Aliases, []string{"hi"}) {
		t.Errorf("aliases = %v, want [hi]", got.Aliases)
	}
}

// A name that a staged alias already claimed is renamed, even under
// --overwrite: --overwrite replaces commands on disk, not entries in the pack.
func TestImportRenamesAPackInternalCollision(t *testing.T) {
	h := newHarness()
	h.packs.read = []catalog.Command{
		{Name: "a", Aliases: []string{"b"}, Exec: "echo a"},
		{Name: "b", Exec: "echo b"},
	}
	h.prompt.answers = []bool{true}

	if err := h.useCase().Import("pack.toml", true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if _, found := h.state.saved[0].Get("b_imported"); !found {
		t.Errorf("names = %v, want b_imported", h.state.saved[0].Names())
	}
}

func TestImportRenamesAConflictingAliasWhenConfirmed(t *testing.T) {
	h := newHarness(catalog.Command{Name: "other", Aliases: []string{"taken"}, Exec: "echo other"})
	h.packs.read = []catalog.Command{{Name: "greet", Aliases: []string{"taken"}, Exec: "echo hi"}}
	h.prompt.answers = []bool{true}

	if err := h.useCase().Import("pack.toml", false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	got, _ := h.state.saved[0].Get("greet")
	if !slices.Equal(got.Aliases, []string{"taken_imported"}) {
		t.Errorf("aliases = %v, want [taken_imported]", got.Aliases)
	}
}

func TestImportOverwriteDropsAShadowingAlias(t *testing.T) {
	h := newHarness(catalog.Command{Name: "other", Aliases: []string{"taken"}, Exec: "echo other"})
	h.packs.read = []catalog.Command{{Name: "greet", Aliases: []string{"taken"}, Exec: "echo hi"}}

	if err := h.useCase().Import("pack.toml", true); err != nil {
		t.Fatalf("Import: %v", err)
	}

	got, _ := h.state.saved[0].Get("greet")
	if len(got.Aliases) != 0 {
		t.Errorf("aliases = %v, want none", got.Aliases)
	}

	if owner, _ := h.state.saved[0].AliasOwner("taken"); owner != "other" {
		t.Errorf("alias taken = %q, want it to stay on other", owner)
	}
}

func TestImportEmptyPackWritesNothing(t *testing.T) {
	h := newHarness()

	if err := h.useCase().Import("pack.toml", false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	if len(h.state.saved) != 0 {
		t.Errorf("saved %v, want nothing", h.state.saved)
	}
}

func TestImportErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("read", func(t *testing.T) {
		h := newHarness()
		h.packs.readErr = boom

		if err := h.useCase().Import("pack.toml", false); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("prompt", func(t *testing.T) {
		h := newHarness(catalog.Command{Name: "greet", Exec: "echo local"})
		h.packs.read = []catalog.Command{{Name: "greet", Exec: "echo imported"}}
		h.prompt.err = boom

		if err := h.useCase().Import("pack.toml", false); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.saved) != 0 {
			t.Error("a failed prompt still saved")
		}
	})

	t.Run("save", func(t *testing.T) {
		h := newHarness()
		h.packs.read = []catalog.Command{{Name: "greet", Exec: "echo imported"}}
		h.state.saveErr = boom

		if err := h.useCase().Import("pack.toml", false); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

func TestListRegistry(t *testing.T) {
	h := newHarness()
	h.reg.packs = []sharecommands.Pack{
		{Name: "tools", Description: "Useful tools", Version: "1.2.0"},
	}

	if err := h.useCase().List(); err != nil {
		t.Fatalf("List: %v", err)
	}

	if got, want := h.output.results, []string{"tools - Useful tools (v1.2.0)"}; !slices.Equal(got, want) {
		t.Errorf("results = %v, want %v", got, want)
	}
}

func TestPullShowsThePreviewAndImportsOnAYes(t *testing.T) {
	h := newHarness()
	h.reg.pack = []catalog.Command{
		{Name: "b", Exec: "echo b"},
		{Name: "a", Exec: "echo a"},
	}
	h.reg.url = "https://example.test/packs/tools.toml"
	h.prompt.answers = []bool{true}

	if err := h.useCase().Pull("tools"); err != nil {
		t.Fatalf("Pull: %v", err)
	}

	text := strings.Join(h.output.infos, "\n")

	for _, want := range []string{
		"Importing pack 'tools' from https://example.test/packs/tools.toml",
		"Pack contains: a, b",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("diagnostics = %q, want %q", text, want)
		}
	}

	if len(h.state.saved) != 1 {
		t.Fatalf("saved %d catalogs, want the imported pack", len(h.state.saved))
	}
}

func TestPullDeclinedWritesNothing(t *testing.T) {
	h := newHarness()
	h.reg.pack = []catalog.Command{{Name: "a", Exec: "echo a"}}
	h.prompt.answers = []bool{false}

	if err := h.useCase().Pull("tools"); err != nil {
		t.Fatalf("Pull: %v", err)
	}

	if len(h.state.saved) != 0 {
		t.Error("a declined pull imported something")
	}
}

func TestPullErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("fetch", func(t *testing.T) {
		h := newHarness()
		h.reg.packErr = boom

		if err := h.useCase().Pull("tools"); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("index", func(t *testing.T) {
		h := newHarness()
		h.reg.indexErr = boom

		if err := h.useCase().List(); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

// A pack's description and params belong to the command: import must not drop
// them, or a shared pack loses the only thing that explains it.
func TestImportKeepsMetadata(t *testing.T) {
	h := newHarness()
	h.packs.read = []catalog.Command{{
		Name:        "greet",
		Exec:        "echo hi",
		Description: "Say hi",
		Params:      []string{"name"},
		Shell:       "pwsh",
		Cwd:         "~/api",
	}}

	if err := h.useCase().Import("pack.toml", false); err != nil {
		t.Fatalf("Import: %v", err)
	}

	got, found := h.state.saved[0].Get("greet")
	if !found || got.Description != "Say hi" || !slices.Equal(got.Params, []string{"name"}) {
		t.Errorf("imported = %+v, want the metadata", got)
	}

	if got.Shell != "pwsh" || got.Cwd != "~/api" {
		t.Errorf("imported = %+v, want the shell and cwd", got)
	}
}
