// Package logger writes crower's two kinds of line to two different streams.
//
// Results -- the rows of `list`, the output of an executed command, a history
// table -- are data a caller may pipe, so they go to the result stream and keep
// their exact bytes. Everything else (prompts, status, warnings, errors) is a
// diagnostic and goes to the error stream. Both are also written to the log
// file.
//
// The split is what makes the binary composable: before it, every line went to
// stdout through fmt.Println, which also appended a space to every result line.
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// Logger writes results to out and diagnostics to err, and both to a log file.
type Logger struct {
	out     io.Writer
	err     io.Writer
	file    *slog.Logger
	logFile *os.File
}

// New builds a Logger that writes results to out and diagnostics to err, with a
// debug-level copy of both in the file at logpath.
//
// An unwritable log path is a returned error rather than a panic: the log is a
// convenience, and a read-only home directory should not turn `crower list`
// into a crash.
func New(logpath string, out io.Writer, errOut io.Writer) (*Logger, error) {
	if out == nil {
		out = os.Stdout
	}

	if errOut == nil {
		errOut = os.Stderr
	}

	if err := rotateLog(logpath); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file %s: %w", logpath, err)
	}

	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	return &Logger{
		out:     out,
		err:     errOut,
		file:    slog.New(fileHandler),
		logFile: logFile,
	}, nil
}

// Out writes a result line: data a caller may pipe.
//
// The bytes are preserved -- no prefix and no trailing space -- and a newline is
// added only when the text does not already end with one, so a command's own
// output is not rewritten on its way to the stream.
func (logger Logger) Out(msg string) {
	writeLine(logger.out, msg)
	logger.file.Info(msg, "stream", "out")
}

// Info writes a diagnostic line: status, progress, a prompt.
func (logger Logger) Info(msg string, data ...any) {
	writeLine(logger.err, joinData(msg, data))
	logger.file.Info(msg, "data", dataToString(data))
}

// Warning writes a diagnostic line the user should read but that does not fail
// the run.
func (logger Logger) Warning(msg string, data ...any) {
	writeLine(logger.err, joinData(msg, data))
	logger.file.Warn(msg, "data", dataToString(data))
}

// Error writes a diagnostic line for a failure.
func (logger Logger) Error(msg string, data ...any) {
	writeLine(logger.err, joinData(msg, data))
	logger.file.Error(msg, "data", dataToString(data))
	logger.stack()
}

// Debug writes a line only the log file receives.
//
// It is deliberately silent on the terminal: a step marker belongs in the log,
// not in the middle of a pipe. This is what `crower execute` uses to record the
// exact shell line it ran.
func (logger Logger) Debug(msg string, data ...any) {
	logger.file.Debug(msg, "data", dataToString(data))
}

// Close releases the log file. A run that ends normally still has the file
// closed by the process exit, but a test that opens hundreds of them does not.
func (logger Logger) Close() error {
	if logger.logFile == nil {
		return nil
	}

	return logger.logFile.Close()
}

// writeLine writes msg followed by exactly one newline.
func writeLine(w io.Writer, msg string) {
	if strings.HasSuffix(msg, "\n") {
		io.WriteString(w, msg)
		return
	}

	io.WriteString(w, msg+"\n")
}

// joinData appends the values to msg, separated by a space.
//
// The values are not key/value pairs: a caller that passes ("Name: ", "name",
// x) is writing a sentence, and the old writer printed the labels as words.
func joinData(msg string, data []any) string {
	if len(data) == 0 {
		return msg
	}

	return strings.TrimSuffix(msg, " ") + " " + dataToString(data)
}

// dataToString writes every value it is given.
func dataToString(data []any) string {
	builder := strings.Builder{}
	for _, d := range data {
		fmt.Fprintf(&builder, "%v ", d)
	}

	return strings.TrimSpace(builder.String())
}

func (logger Logger) stack() {
	if logger.logFile == nil {
		return
	}

	fmt.Fprintln(logger.logFile, "Stacktrace: ")
	for line := range strings.SplitSeq(string(debug.Stack()), "\n") {
		fmt.Fprintln(logger.logFile, line)
	}
}

// rotateLog keeps the previous log next to the current one, so a run does not
// destroy the evidence of the run before it.
func rotateLog(logpath string) error {
	dirpath := filepath.Dir(logpath)
	base := "old_" + filepath.Base(logpath)

	if err := copyFile(logpath, filepath.Join(dirpath, base)); err != nil {
		return fmt.Errorf("rotate log %s: %w", logpath, err)
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}
