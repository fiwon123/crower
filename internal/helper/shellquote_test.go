package helper

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
		{"with quote", `he said "hi"`, `"he said \"hi\""`},
		{"with percent", "%PATH%", `"%%PATH%%"`},
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
