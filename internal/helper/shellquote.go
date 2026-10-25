package helper

import (
	"runtime"
	"strings"
)

// QuoteFor quotes a string for safe use in shell commands on the given OS.
// goos should be runtime.GOOS or "windows"/"linux" for testing.
func QuoteFor(goos, s string) string {
	if goos == "windows" {
		return quoteWindows(s)
	}
	return quotePosix(s)
}

// Quote quotes for the current OS.
func Quote(s string) string {
	return QuoteFor(runtime.GOOS, s)
}

func quotePosix(s string) string {
	if s == "" {
		return "''"
	}
	// Use single quotes; escape single quotes as '\''
	if !strings.Contains(s, "'") {
		return "'" + s + "'"
	}
	// Replace ' with '\'' 
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func quoteWindows(s string) string {
	if s == "" {
		return `""`
	}
	// For cmd.exe, wrap in double quotes and escape internal double quotes
	// Also escape percent signs to avoid variable expansion
	if strings.ContainsAny(s, `"%^&|<>`) {
		s = strings.ReplaceAll(s, "%", "%%")
		s = strings.ReplaceAll(s, `"`, `\"`)
		return `"` + s + `"`
	}
	// Simple case - still quote if contains spaces or special chars
	if strings.ContainsAny(s, " \t\n\r(),;=") {
		s = strings.ReplaceAll(s, "%", "%%")
		s = strings.ReplaceAll(s, `"`, `\"`)
		return `"` + s + `"`
	}
	return s
}
