//go:build fyne

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

// windowLog is the record of one window run.
//
// The window executable is linked as a GUI program, so it has no console: a
// start that fails inside Fyne writes to a logger nobody can see and then exits
// the process. Without this file, that failure is a double-click that does
// nothing (#465).
type windowLog struct {
	path string
	file *os.File
}

// startWindowLog opens <config>/fyne.log and points the standard logger — the
// one Fyne reports through — at it, keeping stderr as a second sink when there
// is one.
//
// The first Fyne error is also handed to showFyneStartupError, because Fyne
// exits the process itself on a window it cannot create: that is the only
// moment a dialog can still be shown.
func startWindowLog(folder string) *windowLog {
	path := filepath.Join(folder, "fyne.log")

	file, err := os.Create(path)
	if err != nil {
		return &windowLog{}
	}

	record := &windowLog{path: path, file: file}

	writer := io.MultiWriter(file, os.Stderr)
	log.SetOutput(&startupErrors{writer: writer, path: path})
	log.Printf("crower-gui %s starting (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)

	return record
}

// clean ends a run that closed normally. The log is only worth keeping when
// something went wrong, so a good run leaves nothing behind.
func (l *windowLog) clean() {
	if l == nil || l.file == nil {
		return
	}

	log.SetOutput(os.Stderr)
	_ = l.file.Close()

	path := l.path
	l.file = nil
	l.path = ""

	_ = os.Remove(path)
}

// keep ends a failed run and leaves the file for the user to send.
func (l *windowLog) keep() {
	if l == nil || l.file == nil {
		return
	}

	log.SetOutput(os.Stderr)
	_ = l.file.Close()
	l.file = nil
}

// recordPanic appends a recovered panic and its stack to the file.
func (l *windowLog) recordPanic(recovered any) {
	if l == nil || l.path == "" {
		return
	}

	file, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "panic: %v\n%s", recovered, debug.Stack())
}

// startupErrors tees the log to the file and, on the first Fyne error, shows a
// dialog. Fyne reports a failed window through this logger and then exits the
// process, so the dialog has to be raised from the write itself.
type startupErrors struct {
	writer io.Writer

	path     string
	reported bool
	mutex    sync.Mutex
}

func (s *startupErrors) Write(data []byte) (int, error) {
	written, err := s.writer.Write(data)

	line := strings.TrimSpace(string(data))
	if line == "" {
		return written, err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.reported {
		return written, err
	}

	// Fyne writes a failed window start as "Fyne error: ...". A dialog for the
	// first one is enough: a second would repeat the same story.
	if !strings.Contains(line, "Fyne error") {
		return written, err
	}

	s.reported = true

	message := "crower-gui could not open its window.\n\n" + line
	if s.path != "" {
		message += "\n\nDetails: " + s.path
	}

	notifyStartupError(message)

	return written, err
}

// notifyStartupError is showFyneStartupError, and a variable so a test can see
// that the first Fyne error raises exactly one notification.
var notifyStartupError = showFyneStartupError
