package helper

import (
	"runtime"

	"github.com/fiwon123/crower/internal/shell"
)

// QuoteFor quotes a string for safe use in shell commands on the given OS.
// goos should be runtime.GOOS or "windows"/"linux" for testing.
//
// The implementation lives in internal/shell so the note writer can quote a
// recorded command line without importing this package, which depends on app.
func QuoteFor(goos, s string) string {
	return shell.QuoteFor(goos, s)
}

// Quote quotes for the current OS.
func Quote(s string) string {
	return QuoteFor(runtime.GOOS, s)
}
