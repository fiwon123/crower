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

// ErrChecksumMismatch reports a download whose bytes do not match the release's
// checksums.txt.
var ErrChecksumMismatch = stderrors.New("checksum mismatch")

// ErrChecksumUnavailable reports that checksums.txt is missing, unreadable, or
// has no usable entry for the asset.
var ErrChecksumUnavailable = stderrors.New("checksum unavailable")

// ErrUnknownShell reports a command that names a shell crower does not know.
// The portable schema allows "sh", "pwsh" and "cmd"; a typo must fail before a
// process starts rather than fall back to a shell the command did not ask for.
var ErrUnknownShell = stderrors.New("unknown shell")

// ErrInvalidWorkDir reports a working directory that does not exist or is not a
// directory. A command that names one must fail before anything runs.
var ErrInvalidWorkDir = stderrors.New("invalid working directory")

// ErrInvalidTimeout reports a timeout that is not a positive duration. A zero
// or negative timeout would kill every step immediately, so it is refused
// rather than obeyed.
var ErrInvalidTimeout = stderrors.New("invalid timeout")

// ErrInvalidEnv reports an environment entry that cannot be handed to a
// process: an empty key, or a key or value carrying "=" or NUL.
var ErrInvalidEnv = stderrors.New("invalid environment entry")

// ErrSecretEnvValue reports a literal value under a key that looks like a
// secret. The portable schema accepts a reference ("${VAR}") there, never the
// secret itself: a pack is a file people commit and share.
var ErrSecretEnvValue = stderrors.New("secret-looking environment value")

// ErrUnsetEnvVar reports a "${VAR}" reference to a variable the process
// environment does not define. Passing an empty value would make the command
// run with a missing credential instead of saying so.
var ErrUnsetEnvVar = stderrors.New("environment variable is not set")

// ErrTimeout reports a step killed because it exceeded its command's timeout.
var ErrTimeout = stderrors.New("step timed out")

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
