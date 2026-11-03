package searchdom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestBrowserCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		content string
		want    string
		wantErr bool
	}{
		{
			"windows URL",
			"windows", "https://example.com",
			`start ' ' 'https://example.com' `, false,
		},
		{
			"windows search",
			"windows", "golang context",
			`start ' ' 'https://duckduckgo.com/?q=golang context' `, false,
		},
		{
			"linux URL",
			"linux", "https://example.com",
			`xdg-open "https://example.com" `, false,
		},
		{
			"linux search",
			"linux", "golang",
			`xdg-open "https://duckduckgo.com/?q=golang" `, false,
		},
		{
			"darwin URL",
			"darwin", "https://example.com",
			`open "https://example.com"`, false,
		},
		{
			"darwin search",
			"darwin", "golang",
			`open "https://duckduckgo.com/?q=golang"`, false,
		},
		{"unsupported platform", "plan9", "golang", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := browserCommand(tt.goos, tt.content)

			if tt.wantErr {
				if !errors.Is(err, crowererrors.ErrUnsupportedOS) {
					t.Fatalf("error = %v, want ErrUnsupportedOS", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("error = %v", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSearchCommands(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		content string
		command func(goos, path, content string) (string, error)
		want    string
		wantErr bool
	}{
		{"files on windows", "windows", "a.txt", searchFileCommand, `dir "." /s /b /a-d | findstr "a.txt"`, false},
		{"files on linux", "linux", "a.txt", searchFileCommand, `find "." -type f -name "a.txt"`, false},
		{"files on darwin", "darwin", "a.txt", searchFileCommand, `find "." -type f -name "a.txt"`, false},
		{"folders on windows", "windows", "a", searchFolderCommand, `dir "." /s /b /ad | findstr "a"`, false},
		{"folders on linux", "linux", "a", searchFolderCommand, `find "." -type d -name "a"`, false},
		{"folders on darwin", "darwin", "a", searchFolderCommand, `find "." -type d -name "a"`, false},
		{"both on windows", "windows", "a", searchFileAndFolderCommand, `dir "." /s /b | findstr "a"`, false},
		{"both on linux", "linux", "a", searchFileAndFolderCommand, `find "." -name "a"`, false},
		{"both on darwin", "darwin", "a", searchFileAndFolderCommand, `find "." -name "a"`, false},
		{"files on an unsupported platform", "plan9", "a.txt", searchFileCommand, "", true},
		{"folders on an unsupported platform", "plan9", "a", searchFolderCommand, "", true},
		{"both on an unsupported platform", "plan9", "a", searchFileAndFolderCommand, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.command(tt.goos, ".", tt.content)

			if tt.wantErr {
				if !errors.Is(err, crowererrors.ErrUnsupportedOS) {
					t.Fatalf("error = %v, want ErrUnsupportedOS", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("error = %v", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}
