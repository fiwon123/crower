package main

import (
	"strings"
	"testing"
)

func TestCheckNames(t *testing.T) {
	tests := []struct {
		name    string
		entries []string
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "plain unix names",
			entries: []string{"crower", "cr", "add_to_path.sh", "README.md"},
			wantOK:  true,
		},
		{
			name:    "plain windows names",
			entries: []string{"crower.exe", "cr.exe", "add_to_path.bat"},
			wantOK:  true,
		},
		{
			name:    "the v0.4.0 regression",
			entries: []string{"crower_linux_amd64", "cr_linux_amd64", "add_to_path.sh"},
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "the v0.4.0 windows regression",
			entries: []string{"crower_windows_amd64.exe", "cr_windows_amd64.exe"},
			wantMsg: "arch-suffixed binary",
		},
		{
			name:    "short name missing",
			entries: []string{"crower", "add_to_path.sh"},
			wantMsg: "missing \"cr\"",
		},
		{
			name:    "long name missing",
			entries: []string{"cr"},
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "empty archive",
			entries: nil,
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "parent traversal",
			entries: []string{"crower", "cr", "../README.md"},
			wantMsg: "unsafe member",
		},
		{
			name:    "absolute member",
			entries: []string{"crower", "cr", "/etc/passwd"},
			wantMsg: "unsafe member",
		},
		{
			name:    "windows parent traversal",
			entries: []string{"crower.exe", "cr.exe", `..\README.md`},
			wantMsg: "unsafe member",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems := checkNames(tt.entries)
			if tt.wantOK {
				if len(problems) != 0 {
					t.Fatalf("checkNames(%v) = %v, want no problems", tt.entries, problems)
				}
				return
			}
			if len(problems) == 0 {
				t.Fatalf("checkNames(%v) passed, want a problem", tt.entries)
			}
			joined := strings.Join(problems, "\n")
			if !strings.Contains(joined, tt.wantMsg) {
				t.Errorf("checkNames(%v) = %q, want it to mention %q", tt.entries, joined, tt.wantMsg)
			}
		})
	}
}

func TestIsUnsafePath(t *testing.T) {
	tests := []struct {
		entry string
		want  bool
	}{
		{"crower", false},
		{"README.md", false},
		{"add_to_path.sh", false},
		{"../README.md", true},
		{"a/../b", true},
		{`..\README.md`, true},
		{"/etc/passwd", true},
		{`C:\Windows\system32`, true},
		{"a..b", false},
	}

	for _, tt := range tests {
		if got := isUnsafePath(tt.entry); got != tt.want {
			t.Errorf("isUnsafePath(%q) = %v, want %v", tt.entry, got, tt.want)
		}
	}
}

func TestIsArchSuffixed(t *testing.T) {
	tests := []struct {
		base string
		want bool
	}{
		{"crower_linux_amd64", true},
		{"cr_linux_amd64", true},
		{"crower_windows_amd64.exe", true},
		{"cr_darwin_arm64", true},
		{"crower", false},
		{"cr", false},
		{"crower.exe", false},
		{"README.md", false},
		{"add_to_path.sh", false},
		{"crowex_linux_amd64", false},
	}

	for _, tt := range tests {
		if got := isArchSuffixed(tt.base); got != tt.want {
			t.Errorf("isArchSuffixed(%q) = %v, want %v", tt.base, got, tt.want)
		}
	}
}
