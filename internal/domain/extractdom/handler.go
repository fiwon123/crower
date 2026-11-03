package extractdom

import (
	stderrors "errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/shell"
)

type Handler struct {
	app   *app.Config
	shell *shell.Runner
}

func NewHandler(app *app.Config) *Handler {

	return &Handler{
		app:   app,
		shell: shell.New(),
	}
}

// Extract expands every archive in paths.
//
// Every archive is attempted; if any fails, the failures are returned together
// so the command exits non-zero. The previous version logged a failure and
// carried on, which a script read as success.
func (h *Handler) Extract(paths []string, output string) error {
	var failures []error

	for _, f := range paths {
		// The destination belongs to this archive: an omitted --out means
		// "next to this archive", not "next to the first archive". Assigning
		// it once outside the loop sent every later archive into the first
		// one's directory.
		outDir := output
		if outDir == "" {
			outDir = filepath.Dir(f)
		}

		base := filepath.Base(f)
		split := strings.Split(base, ".")
		ext := ""

		if len(split) > 0 {
			ext = split[len(split)-1]
		}

		if ext == "" {
			failures = append(failures, fmt.Errorf("%s: no file extension to choose an extractor", f))
			continue
		}

		out, err := h.performExtract(ext, f, outDir)
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", f, err))
			continue
		}

		h.app.Logger.Out(strings.TrimRight(string(out), "\n"))
	}

	return stderrors.Join(failures...)
}

func (h *Handler) performExtract(ext string, filePath string, outDir string) (string, error) {
	commandLine, err := extractCommand(runtime.GOOS, ext, filePath, outDir)
	if err != nil {
		return "", err
	}

	return h.shell.Exec(commandLine)
}
