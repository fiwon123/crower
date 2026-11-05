package cmd

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/application/extractcontent"
	"github.com/fiwon123/crower/internal/application/opencontent"
	"github.com/fiwon123/crower/internal/application/searchcontent"
	"github.com/fiwon123/crower/internal/application/transfercontent"
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/fiwon123/crower/internal/session"
)

// fyneTools runs the file operations from the window through the same use
// cases the CLI commands run, and returns what to show instead of printing it.
//
// It has no Fyne import, so the adapter is covered by the normal cmd tests.
type fyneTools struct {
	cfg *session.Session
}

func newFyneTools(cfg *session.Session) fyneTools {
	return fyneTools{cfg: cfg}
}

// Copy copies every source into destination.
func (t fyneTools) Copy(sources []string, destination string) (string, timeline.Log, error) {
	return t.transfer(operation.Copy, sources, destination)
}

// Move moves every source into destination.
func (t fyneTools) Move(sources []string, destination string) (string, timeline.Log, error) {
	return t.transfer(operation.Move, sources, destination)
}

func (t fyneTools) transfer(op operation.Operation, sources []string, destination string) (string, timeline.Log, error) {
	args := make([]string, 0, len(sources)+1)
	args = append(args, sources...)
	args = append(args, destination)

	session := newTransferSession(t.cfg)

	if err := transfercontent.New(session, session, session, session).Run(
		transfercontent.Request{Op: op, Args: args},
	); err != nil {
		return "", t.cfg.History, err
	}

	verb := "Copied"
	if op == operation.Move {
		verb = "Moved"
	}

	return fmt.Sprintf("%s %d item(s) to %s", verb, len(sources), destination), t.cfg.History, nil
}

// Search looks for files, folders, both, or opens a web search.
func (t fyneTools) Search(mode desktop.SearchMode, term, path string) (string, timeline.Log, error) {
	request := searchcontent.Request{Args: []string{term}}

	switch mode {
	case desktop.SearchFiles:
		request.Mode = searchcontent.ModeFile
	case desktop.SearchFolders:
		request.Mode = searchcontent.ModeFolder
	case desktop.SearchBoth:
		request.Mode = searchcontent.ModeBoth
	case desktop.SearchWeb:
		request.Mode = searchcontent.ModeBrowser
	default:
		return "", t.cfg.History, fmt.Errorf("unknown search mode %q", mode)
	}

	if request.Mode != searchcontent.ModeBrowser && path != "" {
		request.Args = append(request.Args, path)
	}

	capture := &captureOutput{}
	search := captureSearchSession{searchSession: newSearchSession(t.cfg), output: capture}

	if err := searchcontent.New(search, search, search, search).Run(request); err != nil {
		return "", t.cfg.History, err
	}

	if request.Mode == searchcontent.ModeBrowser {
		return "Opened a web search for " + term, t.cfg.History, nil
	}

	if text := capture.String(); text != "" {
		return text, t.cfg.History, nil
	}

	return "No matches.", t.cfg.History, nil
}

// Extract expands archives, beside each one unless outDir is set.
func (t fyneTools) Extract(archives []string, outDir string) (string, timeline.Log, error) {
	capture := &captureOutput{}
	extract := captureExtractSession{extractSession: newExtractSession(t.cfg), output: capture}

	if err := extractcontent.New(extract, extract, extract, extract).Run(
		extractcontent.Request{Args: archives, OutDir: outDir},
	); err != nil {
		return "", t.cfg.History, err
	}

	if text := capture.String(); text != "" {
		return text, t.cfg.History, nil
	}

	return fmt.Sprintf("Extracted %d archive(s).", len(archives)), t.cfg.History, nil
}

// Open opens paths with the platform's default application.
func (t fyneTools) Open(paths []string) (string, timeline.Log, error) {
	capture := &captureOutput{}
	open := captureOpenSession{openSession: newOpenSession(t.cfg), output: capture}

	if err := opencontent.New(open, open, open, open, open).Run(
		opencontent.Request{Mode: opencontent.ModePaths, Args: paths},
	); err != nil {
		return "", t.cfg.History, err
	}

	if text := capture.String(); text != "" {
		return text, t.cfg.History, nil
	}

	return "Opened " + strings.Join(paths, ", "), t.cfg.History, nil
}

// captureOutput collects what a use case would print, so the window can show
// it in its result pane instead of a terminal that may not exist.
//
// One type serves the search, extract and open output ports: they differ only
// in which of these three methods they declare.
type captureOutput struct {
	lines []string
}

func (c *captureOutput) add(text string) {
	text = strings.TrimRight(text, "\n")
	if strings.TrimSpace(text) == "" {
		return
	}

	c.lines = append(c.lines, text)
}

// Info implements the open use case's diagnostics stream.
func (c *captureOutput) Info(message string, data ...any) {
	var line strings.Builder
	line.WriteString(message)

	for _, value := range data {
		line.WriteString(fmt.Sprint(value))
	}

	c.add(line.String())
}

// Result implements the extract and open result stream.
func (c *captureOutput) Result(text string) {
	c.add(text)
}

// Emit implements the search result stream. The structured payload is not used
// here: the window shows the text the search command printed.
func (c *captureOutput) Emit(text string, _ any) error {
	c.add(text)

	return nil
}

// String is everything captured, one entry per line.
func (c *captureOutput) String() string {
	return strings.Join(c.lines, "\n")
}

// captureSearchSession redirects the search result into the window instead of
// the process streams; everything else comes from the CLI's own session wiring.
type captureSearchSession struct {
	searchSession
	output *captureOutput
}

func (s captureSearchSession) Emit(text string, payload any) error {
	return s.output.Emit(text, payload)
}

// captureExtractSession redirects the extractor's output into the window.
type captureExtractSession struct {
	extractSession
	output *captureOutput
}

func (s captureExtractSession) Result(text string) {
	s.output.Result(text)
}

// captureOpenSession redirects the opener's diagnostics and result into the
// window.
type captureOpenSession struct {
	openSession
	output *captureOutput
}

func (s captureOpenSession) Info(message string, data ...any) {
	s.output.Info(message, data...)
}

func (s captureOpenSession) Result(text string) {
	s.output.Result(text)
}
