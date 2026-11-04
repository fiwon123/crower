package errors

import (
	stderrors "errors"
	"fmt"
	"testing"
)

// The sentinels are the contract callers test with errors.Is. Pin the text and
// make sure each one is its own value: two names aliasing one error would make
// errors.Is answer for the wrong kind.
func TestSentinels(t *testing.T) {
	sentinels := []struct {
		name    string
		err     error
		message string
	}{
		{"unsupported OS", ErrUnsupportedOS, "unsupported operating system"},
		{"invalid pack name", ErrInvalidPackName, "invalid pack name"},
		{"checksum mismatch", ErrChecksumMismatch, "checksum mismatch"},
		{"checksum unavailable", ErrChecksumUnavailable, "checksum unavailable"},
		{"unknown shell", ErrUnknownShell, "unknown shell"},
		{"invalid work dir", ErrInvalidWorkDir, "invalid working directory"},
		{"invalid timeout", ErrInvalidTimeout, "invalid timeout"},
		{"invalid env", ErrInvalidEnv, "invalid environment entry"},
		{"secret env", ErrSecretEnvValue, "secret-looking environment value"},
		{"unset env var", ErrUnsetEnvVar, "environment variable is not set"},
		{"timeout", ErrTimeout, "step timed out"},
	}

	seen := make(map[error]string, len(sentinels))

	for _, tt := range sentinels {
		if tt.err == nil {
			t.Fatalf("%s: sentinel is nil", tt.name)
		}

		if got := tt.err.Error(); got != tt.message {
			t.Errorf("%s: message = %q, want %q", tt.name, got, tt.message)
		}

		if !stderrors.Is(fmt.Errorf("%w: detail", tt.err), tt.err) {
			t.Errorf("%s: errors.Is did not find the wrapped sentinel", tt.name)
		}

		if previous, aliased := seen[tt.err]; aliased {
			t.Errorf("%s aliases %s", tt.name, previous)
		}

		seen[tt.err] = tt.name
	}
}

// The constructors have no sentinel, so the message is the contract: it is
// what the user reads on stderr.
func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"not enough args", NotEnoughArgs("path"), "missing arguments: path"},
		{"too many args", TooManyArgs("one name"), "too many arguments: expected one name"},
		{"command not found", CommandNotFound("nope"), `command not found: "nope"`},
		{"file and folder flags", FileAndFolderFlagsNotUsed(), "--file and --folder cannot be used together"},
		{"no path", NotFileAndOutputPath(), "needs a file or folder path"},
		{"empty paths", EmptyPaths(), "no paths given"},
		{"no user home", NotUserHome(), "cannot determine the user home directory"},
		{"input required", InputRequired("name"), `input required for "name" but --no-input is set`},
		{"no input", NoInputAvailable(), "no input available: standard input is closed"},
		{"invalid rows", InvalidRows(), "invalid row"},
		{"invalid input", InvalidInput(), "invalid input"},
		{"empty input", EmptyInput(), "input is empty"},
		{"letters and numbers", OnlyLettersAndNumbers(), "only letters and numbers are allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatal("constructor returned nil")
			}

			if got := tt.err.Error(); got != tt.want {
				t.Errorf("message = %q, want %q", got, tt.want)
			}
		})
	}
}
