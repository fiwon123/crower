package shell

import (
	"runtime"
	"strings"
)

// QuoteFor quotes a string for safe use in shell commands on the given OS.
// goos should be runtime.GOOS or "windows"/"linux" for testing.
//
// It lives here rather than in internal/helper because the note writer uses it
// to record a replayable command line, and helper imports app, which imports
// the note package.
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

// quoteWindows quotes one argument in a line handed to cmd /c.
//
// cmd has no backslash escape. An embedded quote is written as "" — cmd
// toggles its quoting twice, so the state around it does not change, and the
// program's argument parser reads the pair as one literal quote. Percent signs
// are left alone: "%%" is a batch-file escape that a command line does not
// collapse, so doubling corrupts every path containing "%". A "%VAR%" still
// expands, which no command-line quoting can prevent.
//
// A trailing run of backslashes is doubled when the value is quoted: a single
// "\" before the closing quote escapes it, and the argument swallows whatever
// follows. A backslash immediately before an embedded quote cannot be carried
// exactly: the run is doubled so the quote state stays balanced, which is the
// safe choice even though the value is not preserved byte for byte.
func quoteWindows(s string) string {
	if s == "" {
		return `""`
	}

	if !needsWindowsQuote(s) {
		return s
	}

	var quoted strings.Builder
	quoted.WriteByte('"')

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			quoted.WriteString(`""`)

		case '\\':
			end := i
			for end < len(s) && s[end] == '\\' {
				end++
			}

			run := s[i:end]
			quoted.WriteString(run)

			// Both a trailing run and a run before an embedded quote must be
			// even: an odd count turns the next quote into a literal escape
			// and leaves the quoting state flipped, so a later "&" would run.
			// The cost is that a backslash directly before a literal quote is
			// not preserved; no cmd line can carry both.
			if end == len(s) || s[end] == '"' {
				quoted.WriteString(run)
			}

			i = end - 1

		default:
			quoted.WriteByte(s[i])
		}
	}

	quoted.WriteByte('"')

	return quoted.String()
}

// needsWindowsQuote reports whether cmd would read s as more than one word, or
// as a command separator, unless it is quoted.
func needsWindowsQuote(s string) bool {
	return strings.ContainsAny(s, "\" \t\n\r&|<>^(),;=")
}
