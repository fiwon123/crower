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
// The target program parses argv with the Windows rules (a \" is a literal
// quote), so an embedded quote is written as \". cmd has no backslash escape,
// so that same \" toggles cmd's quoting state; while cmd is outside a quote,
// the characters that would otherwise separate commands (space, &, |, <, >,
// ^, parentheses) are prefixed with ^. The child sees the original value, and
// cmd does not run the metacharacters.
//
// Percent signs are left alone: "%%" is a batch-file escape that a command line
// does not collapse, so doubling corrupts every path containing "%". A "%VAR%"
// still expands, which no command-line quoting can prevent.
//
// A backslash before a quote is doubled the way the Windows argv parser
// expects: 2n+1 backslashes before a literal quote keep n backslashes, and a
// trailing run is doubled before the closing quote.
func quoteWindows(s string) string {
	if s == "" {
		return `""`
	}

	if !needsWindowsQuote(s) {
		return s
	}

	var quoted strings.Builder
	quoted.WriteByte('"')

	// cmd toggles its quote state on every unescaped ". The opening quote
	// starts inside one; an embedded \" toggles it.
	cmdQuoted := true
	backslashes := 0

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {
		case '\\':
			backslashes++

		case '"':
			quoted.WriteString(strings.Repeat(`\`, 2*backslashes+1))
			quoted.WriteByte('"')
			backslashes = 0
			cmdQuoted = !cmdQuoted

		default:
			quoted.WriteString(strings.Repeat(`\`, backslashes))
			backslashes = 0

			if !cmdQuoted && needsCmdEscape(c) {
				quoted.WriteByte('^')
			}

			quoted.WriteByte(c)
		}
	}

	// A trailing run must be even, or the closing quote would be escaped for
	// the child's argv parser.
	quoted.WriteString(strings.Repeat(`\`, 2*backslashes))
	quoted.WriteByte('"')

	return quoted.String()
}

// needsCmdEscape reports whether c is a cmd metacharacter (or whitespace) that
// must be prefixed with ^ while cmd is outside a quote.
func needsCmdEscape(c byte) bool {
	switch c {
	case ' ', '\t', '&', '|', '<', '>', '^', '(', ')':
		return true
	default:
		return false
	}
}

// needsWindowsQuote reports whether cmd would read s as more than one word, or
// as a command separator, unless it is quoted.
func needsWindowsQuote(s string) bool {
	return strings.ContainsAny(s, "\" \t\n\r&|<>^(),;=")
}
