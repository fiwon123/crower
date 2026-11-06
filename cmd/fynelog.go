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
	"time"
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

// openGLHint explains the one failure that has nothing to do with the download:
// Fyne needs a driver that provides OpenGL 2.1, and Remote Desktop and the basic
// display driver do not.
const openGLHint = "The display driver does not provide OpenGL 2.1. That is what " +
	"Remote Desktop and the basic display driver report; a local session with the " +
	"graphics driver installed is needed. `crower gui` offers the same tools in a terminal."

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
	fatal := isFatalWindowError(line)

	// Fyne writes its reason on the next line, so the dialog waits for it
	// instead of reporting a sentence that ends at the colon.
	go s.report(fatal)

	return written, err
}

// report collects the lines Fyne just wrote, shows them, and ends a run whose
// window cannot exist.
func (s *startupErrors) report(fatal bool) {
	time.Sleep(250 * time.Millisecond)

	details := s.details()
	if details == "" {
		details = "Fyne reported an error while opening the window."
	}

	message := "crower-gui could not open its window.\n\n" + details
	if s.path != "" {
		message += "\n\nLog: " + s.path
	}

	if fatal {
		message += "\n\n" + openGLHint
	}

	notifyStartupError(message)

	// Without a window the process would sit in the task manager doing
	// nothing, which is what "nothing happens" looks like to the user.
	if fatal {
		os.Exit(1)
	}
}

// details is the log written so far, from the first Fyne error on.
func (s *startupErrors) details() string {
	if s.path == "" {
		return ""
	}

	body, err := os.ReadFile(s.path)
	if err != nil {
		return ""
	}

	if index := strings.Index(string(body), "Fyne error"); index >= 0 {
		return strings.TrimSpace(string(body[index:]))
	}

	return strings.TrimSpace(string(body))
}

// isFatalWindowError reports whether Fyne could not create a window at all.
func isFatalWindowError(line string) bool {
	return strings.Contains(line, "window creation error") ||
		strings.Contains(line, "initialise GLFW")
}

// notifyStartupError is showFyneStartupError, and a variable so a test can see
// that the first Fyne error raises exactly one notification.
var notifyStartupError = showFyneStartupError
