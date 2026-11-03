package shell

import "testing"

func TestQuoteForPosix(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "''"},
		{"simple", "hello", "'hello'"},
		{"with space", "hello world", "'hello world'"},
		{"with single quote", "it's", "'it'\\''s'"},
		{"with double quote", `he said "hi"`, `'he said "hi"'`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuoteFor("linux", tt.in)
			if got != tt.want {
				t.Errorf("QuoteFor(linux, %q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestQuoteForWindows(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", `""`},
		{"simple no special", "hello", "hello"},
		{"with space", "hello world", `"hello world"`},
		{"with quote", `he said "hi"`, `"he said ""hi"""`},
		{"only a quote", `"`, `""""`},
		{"quote before a separator", `a" & calc`, `"a"" & calc"`},
		{"backslash before a quote", `x\"y`, `"x\\""y"`},
		{"with percent", "%PATH%", "%PATH%"},
		{"percent inside a quoted value", "%PATH% x", `"%PATH% x"`},
		{"with ampersand", "a & calc", `"a & calc"`},
		{"trailing backslash when quoted", `C:\My Files\`, `"C:\My Files\\"`},
		{"two trailing backslashes when quoted", `a b\\`, `"a b\\\\"`},
		{"trailing backslash unquoted", `C:\out\`, `C:\out\`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuoteFor("windows", tt.in)
			if got != tt.want {
				t.Errorf("QuoteFor(windows, %q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
