// Package errors holds the errors crower returns to its command layer.
//
// It deliberately does not import internal/app. An error is a value, and a
// package that builds values must not also own the terminal: the old helpers
// here printed through app.Logger and returned nothing, so no caller could
// propagate a failure and every command exited 0. Callers now decide where a
// message is printed, which is exactly what lets cmd report a non-zero status.
package errors

import (
	stderrors "errors"
	"fmt"
)

// ErrUnsupportedOS reports an operation that has no implementation for the
// running platform. It is an error rather than a silent no-op: an operation
// that did nothing must not look like one that succeeded.
var ErrUnsupportedOS = stderrors.New("unsupported operating system")

// ErrInvalidPackName reports a pack name that is not a registry slug. The name
// is interpolated into a URL and used to derive a file, so a path or a URL
// metacharacter must be rejected before either happens.
var ErrInvalidPackName = stderrors.New("invalid pack name")

// NotEnoughArgs reports a command that was run without the arguments it needs.
// what names the missing input, so the message says what to add.
func NotEnoughArgs(what string) error {
	return fmt.Errorf("missing arguments: %s", what)
}

// TooManyArgs reports a command given more arguments than it accepts. what
// describes what the command expects, so the message says what to remove.
func TooManyArgs(what string) error {
	return fmt.Errorf("too many arguments: expected %s", what)
}

// CommandNotFound reports a name or alias that matches no managed command.
func CommandNotFound(key string) error {
	return fmt.Errorf("command not found: %q", key)
}

// FileAndFolderFlagsNotUsed reports --file and --folder used together, which
// ask for two different kinds of path.
func FileAndFolderFlagsNotUsed() error {
	return fmt.Errorf("--file and --folder cannot be used together")
}

// NotFileAndOutputPath reports a file operation that has no path to work on.
func NotFileAndOutputPath() error {
	return fmt.Errorf("needs a file or folder path")
}

// EmptyPaths reports a set of paths that matched nothing.
func EmptyPaths() error {
	return fmt.Errorf("no paths given")
}

// NotUserHome reports that the home directory could not be resolved, so the
// default configuration path cannot be built.
func NotUserHome() error {
	return fmt.Errorf("cannot determine the user home directory")
}
