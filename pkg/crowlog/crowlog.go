// Package wrapper using zap.
// Provides a easy way to print like fmt package.
package crowlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"log/slog"
)

// The Crowlog wrapper to access logger provider.
type Logger struct {
	terminalLogger *slog.Logger
	fullTerminal   bool
	fileLogger     *slog.Logger
	logFile        *os.File
}

// Create a new LoggerInfo pointer.
func New(logpath string, verbose bool) *Logger {

	terminalLevel := slog.LevelInfo
	if verbose {
		terminalLevel = slog.LevelDebug
	}

	fileLevel := slog.LevelDebug

	dirpath := filepath.Dir(logpath)
	base := filepath.Base(logpath)
	base = "old_" + base
	err := copyFile(logpath, filepath.Join(dirpath, base))
	if err != nil {
		panic("Failed to copy to old log file: " + err.Error())
	}

	logFile, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	termHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: terminalLevel,
	})

	// File handler
	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level:     fileLevel,
		AddSource: true,
	})

	newLogger := &Logger{
		terminalLogger: slog.New(termHandler),
		fullTerminal:   false,
		fileLogger:     slog.New(fileHandler),
		logFile:        logFile}

	return newLogger
}

// Print info message and data of any type.
// data can be string, int, slices, etc.
func (logger Logger) Info(msg string, data ...any) {
	if logger.fullTerminal {
		logger.terminalLogger.Info(msg, data...)
	} else {
		printTerminal(msg, data)
	}

	logger.fileLogger.Info(msg, data...)
}

// Print error message and data of any type.
// data can be string, int, slices, etc.
func (logger Logger) Error(msg string, data ...any) {
	if logger.fullTerminal {
		logger.terminalLogger.Error(msg, data...)
	} else {
		printTerminal(msg, data)
	}

	logger.fileLogger.Error(msg, data...)
	stack(logger.logFile)
}

// Print warning message and data of any type.
// data can be string, int, slices, etc.
func (logger Logger) Warning(msg string, data ...any) {
	if logger.fullTerminal {
		logger.terminalLogger.Warn(msg, data...)
	} else {
		printTerminal(msg, data)
	}

	logger.fileLogger.Warn(msg, data...)
}

// Print debug message and data of any type.
// data can be string, int, slices, etc.
func (logger Logger) Debug(msg string, data ...any) {
	if logger.fullTerminal {
		logger.terminalLogger.Debug(msg, data...)
	} else {
		printTerminal(msg, data)
	}

	logger.fileLogger.Debug(msg, data...)
}

func printTerminal(msg string, data []any) {
	dataBuilder := strings.Builder{}
	for i := 1; i < len(data); i += 2 {
		d := data[i]
		fmt.Fprintf(&dataBuilder, "%v ", d)
	}

	fmt.Println(msg, dataBuilder.String())
}

func stack(logFile *os.File) {
	fmt.Fprintln(logFile, "Stacktrace: ")
	stack := strings.SplitSeq(string(debug.Stack()), "\n")
	for line := range stack {
		fmt.Fprintln(logFile, line)
	}
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
