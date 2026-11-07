//go:build fyne

package fyneui

import (
	"context"
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// toolsTab is the file operations the CLI offers, reachable without leaving the
// window: transfer, search, extract and open. Each card owns its result square,
// which stays hidden until that card runs something.
func (u *ui) toolsTab() fyne.CanvasObject {
	u.transferResult = newToolResult()
	u.searchResult = newToolResult()
	u.extractResult = newToolResult()
	u.openResult = newToolResult()

	return container.NewVScroll(page(container.NewVBox(
		u.transferCard(),
		u.searchCard(),
		u.extractCard(),
		u.openCard(),
	)))
}

// browseButton is the "pick it instead of typing it" control beside a path
// field. Fyne's dialog is single-select, so each press adds one path.
func browseButton(pick func()) *widget.Button {
	button := widget.NewButtonWithIcon("Browse…", theme.FolderOpenIcon(), pick)
	button.Importance = widget.LowImportance

	return button
}

// fieldWithBrowse puts the button at the end of the path field.
func fieldWithBrowse(entry *widget.Entry, pick func()) fyne.CanvasObject {
	return container.NewBorder(nil, nil, nil, browseButton(pick), entry)
}

func (u *ui) transferCard() fyne.CanvasObject {
	u.sourcesEntry = newScrollFreeEntry()
	u.sourcesEntry.SetPlaceHolder("Paths, comma-separated")

	destination := newScrollFreeEntry()
	destination.SetPlaceHolder("Destination folder")

	copyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		paths := splitPaths(u.sourcesEntry.Text)
		u.runTool(u.transferResult, "Copying…", true, func(ctx context.Context, report func(desktop.Progress)) (string, timeline.Log, error) {
			return u.transfer(ctx, report, true, paths, destination.Text)
		})
	})
	copyButton.Importance = widget.HighImportance

	moveButton := widget.NewButtonWithIcon("Move", theme.ContentCutIcon(), func() {
		paths := splitPaths(u.sourcesEntry.Text)
		u.runTool(u.transferResult, "Moving…", true, func(ctx context.Context, report func(desktop.Progress)) (string, timeline.Log, error) {
			return u.transfer(ctx, report, false, paths, destination.Text)
		})
	})
	moveButton.Importance = widget.HighImportance

	form := newFormGrid([]formRow{
		{"Sources", fieldWithBrowse(u.sourcesEntry, func() { u.browseFiles(u.sourcesEntry) })},
		{"Destination", fieldWithBrowse(destination, func() { u.browseFolder(destination) })},
	})

	return newCard(
		"Copy or move",
		"One or more paths and the destination folder.",
		container.NewVBox(form, actionRow(copyButton, moveButton), u.transferResult.object()),
	)
}

var errPathAndDestination = errors.New("enter at least one source path and a destination folder")

// transfer is the call without the goroutine, so a test can drive it.
func (u *ui) transfer(ctx context.Context, report func(desktop.Progress), copying bool, sources []string, destinationText string) (string, timeline.Log, error) {
	destination := strings.TrimSpace(destinationText)
	if len(sources) == 0 || destination == "" {
		return "", timeline.Log{}, errPathAndDestination
	}

	if copying {
		return u.tools.Copy(ctx, report, sources, destination)
	}

	return u.tools.Move(ctx, report, sources, destination)
}

func (u *ui) searchCard() fyne.CanvasObject {
	term := newScrollFreeEntry()
	term.SetPlaceHolder("Search term")

	path := newScrollFreeEntry()
	path.SetPlaceHolder("Folder (default: .)")

	mode := widget.NewSelect([]string{"Files and folders", "Files", "Folders", "Web"}, nil)
	mode.SetSelected("Files and folders")

	searchButton := widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {
		selected := searchMode(mode.Selected)
		u.runTool(u.searchResult, "Searching…", false, func(context.Context, func(desktop.Progress)) (string, timeline.Log, error) {
			return u.search(selected, term.Text, path.Text)
		})
	})
	searchButton.Importance = widget.HighImportance

	form := newFormGrid([]formRow{
		{"Term", term},
		{"Path", fieldWithBrowse(path, func() { u.browseFolder(path) })},
		{"Look in", mode},
	})

	return newCard(
		"Search",
		"File and folder names, or a web search.",
		container.NewVBox(form, actionRow(searchButton), u.searchResult.object()),
	)
}

func searchMode(label string) desktop.SearchMode {
	switch label {
	case "Files":
		return desktop.SearchFiles
	case "Folders":
		return desktop.SearchFolders
	case "Web":
		return desktop.SearchWeb
	default:
		return desktop.SearchBoth
	}
}

// search is the call without the goroutine.
func (u *ui) search(mode desktop.SearchMode, term, pathText string) (string, timeline.Log, error) {
	term = strings.TrimSpace(term)
	if term == "" {
		return "", timeline.Log{}, errors.New("enter a search term")
	}

	path := strings.TrimSpace(pathText)
	if mode == desktop.SearchWeb {
		path = ""
	} else if path == "" {
		path = "."
	}

	return u.tools.Search(mode, term, path)
}

func (u *ui) extractCard() fyne.CanvasObject {
	archives := newScrollFreeEntry()
	archives.SetPlaceHolder("Archive paths, comma-separated")

	outDir := newScrollFreeEntry()
	outDir.SetPlaceHolder("Destination folder (default: beside each archive)")

	extractButton := widget.NewButtonWithIcon("Extract", theme.FileIcon(), func() {
		paths := splitPaths(archives.Text)
		u.runTool(u.extractResult, "Extracting…", true, func(ctx context.Context, report func(desktop.Progress)) (string, timeline.Log, error) {
			return u.extract(ctx, report, paths, outDir.Text)
		})
	})
	extractButton.Importance = widget.HighImportance

	form := newFormGrid([]formRow{
		{"Archives", fieldWithBrowse(archives, func() { u.browseFiles(archives) })},
		{"Into", fieldWithBrowse(outDir, func() { u.browseFolder(outDir) })},
	})

	return newCard(
		"Extract",
		"Zip, tar, tar.gz, 7z and rar archives.",
		container.NewVBox(form, actionRow(extractButton), u.extractResult.object()),
	)
}

// extract is the call without the goroutine.
func (u *ui) extract(ctx context.Context, report func(desktop.Progress), archives []string, outDir string) (string, timeline.Log, error) {
	if len(archives) == 0 {
		return "", timeline.Log{}, errors.New("enter at least one archive path")
	}

	return u.tools.Extract(ctx, report, archives, strings.TrimSpace(outDir))
}

func (u *ui) openCard() fyne.CanvasObject {
	paths := newScrollFreeEntry()
	paths.SetPlaceHolder("Paths, comma-separated")

	openButton := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		u.runTool(u.openResult, "Opening…", false, func(context.Context, func(desktop.Progress)) (string, timeline.Log, error) {
			return u.open(splitPaths(paths.Text))
		})
	})
	openButton.Importance = widget.HighImportance

	return newCard(
		"Open",
		"Open paths with the platform's default application.",
		container.NewVBox(
			fieldWithBrowse(paths, func() { u.browseFiles(paths) }),
			actionRow(openButton),
			u.openResult.object(),
		),
	)
}

// open is the call without the goroutine.
func (u *ui) open(paths []string) (string, timeline.Log, error) {
	if len(paths) == 0 {
		return "", timeline.Log{}, errors.New("enter at least one path")
	}

	return u.tools.Open(paths)
}

// runTool runs the call off the UI thread and reports on it in that card's own
// result area, so a slow search or transfer does not freeze the window.
//
// A cancellable run gets a Cancel button and a progress line; the context stops
// the batch between items, because the platform command already running cannot
// be interrupted.
func (u *ui) runTool(result *toolResult, status string, cancellable bool, action func(ctx context.Context, report func(desktop.Progress)) (string, timeline.Log, error)) {
	if result == nil {
		return
	}

	ctx := context.Background()
	cancel := func() {}

	if cancellable {
		cancellableCtx, stop := context.WithCancel(ctx)
		ctx, cancel = cancellableCtx, stop
		result.start(status, cancel)
	} else {
		result.start(status, nil)
	}

	go func() {
		message, history, err := action(ctx, result.progress)

		fyne.Do(func() {
			cancel()
			result.stop()
			u.reportTool(result, message, history, err)
		})
	}()
}

// reportTool shows the outcome in one card's area and, on success, puts the
// history the operation recorded into the model so the History tab lists it.
func (u *ui) reportTool(result *toolResult, message string, history timeline.Log, err error) {
	if result == nil {
		return
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			result.reportWith("Stopped: "+err.Error(), theme.ColorNameWarning)
			return
		}

		result.reportWith("Error: "+err.Error(), theme.ColorNameError)
		return
	}

	u.model.ReplaceHistory(history)
	u.refreshHistory()

	if strings.TrimSpace(message) == "" {
		message = "Done."
	}

	result.report(message)
}

// splitPaths reads a comma- or newline-separated list, so a pasted command line
// works without editing and a Browse pick can be appended to an earlier one.
func splitPaths(text string) []string {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == ','
	})

	paths := make([]string, 0, len(fields))

	for _, field := range fields {
		if path := strings.TrimSpace(field); path != "" {
			paths = append(paths, path)
		}
	}

	return paths
}
