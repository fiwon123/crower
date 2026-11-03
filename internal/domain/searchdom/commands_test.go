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
			`start "" https://example.com`, false,
		},
		{
			"windows search",
			"windows", "golang context",
			`start "" "https://duckduckgo.com/?q=golang+context"`, false,
		},
		{
			"windows search with a separator",
			"windows", "a & b",
			`start "" "https://duckduckgo.com/?q=a+%26+b"`, false,
		},
		{
			"windows search with a quote",
			"windows", `a" & b`,
			`start "" "https://duckduckgo.com/?q=a%22+%26+b"`, false,
		},
		{
			"windows URL with a query",
			"windows", "https://example.com/?a=1&b=2",
			`start "" "https://example.com/?a=1&b=2"`, false,
		},
		{
			"linux URL",
			"linux", "https://example.com",
			`xdg-open 'https://example.com'`, false,
		},
		{
			"linux search",
			"linux", "golang",
			`xdg-open 'https://duckduckgo.com/?q=golang'`, false,
		},
		{
			"darwin URL",
			"darwin", "https://example.com",
			`open 'https://example.com'`, false,
		},
		{
			"darwin search",
			"darwin", "golang",
			`open 'https://duckduckgo.com/?q=golang'`, false,
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
		{"files on windows", "windows", "a.txt", searchFileCommand, `dir . /s /b /a-d | findstr a.txt`, false},
		{"files on windows with specials", "windows", `a" & calc`, searchFileCommand, `dir . /s /b /a-d | findstr "a"" & calc"`, false},
		{"files on linux", "linux", "a.txt", searchFileCommand, `find '.' -type f -name 'a.txt'`, false},
		{"files on darwin", "darwin", "a.txt", searchFileCommand, `find '.' -type f -name 'a.txt'`, false},
		{"folders on windows", "windows", "a", searchFolderCommand, `dir . /s /b /ad | findstr a`, false},
		{"folders on linux", "linux", "a", searchFolderCommand, `find '.' -type d -name 'a'`, false},
		{"folders on darwin", "darwin", "a", searchFolderCommand, `find '.' -type d -name 'a'`, false},
		{"both on windows", "windows", "a", searchFileAndFolderCommand, `dir . /s /b | findstr a`, false},
		{"both on linux", "linux", "a", searchFileAndFolderCommand, `find '.' -name 'a'`, false},
		{"both on darwin", "darwin", "a", searchFileAndFolderCommand, `find '.' -name 'a'`, false},
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

// The path is quoted too, not only the term: a folder with a space is the
// common case.
func TestSearchCommandQuotesThePath(t *testing.T) {
	tests := []struct {
		name string
		goos string
		path string
		want string
	}{
		{"windows", "windows", `C:\My Files`, `dir "C:\My Files" /s /b /a-d | findstr a.txt`},
		{"windows with an ampersand", "windows", `C:\a&b`, `dir "C:\a&b" /s /b /a-d | findstr a.txt`},
		{"linux", "linux", "/tmp/my files", `find '/tmp/my files' -type f -name 'a.txt'`},
		{"darwin", "darwin", "/tmp/my files", `find '/tmp/my files' -type f -name 'a.txt'`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchFileCommand(tt.goos, tt.path, "a.txt")
			if err != nil {
				t.Fatalf("error = %v", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}

// POSIX quoting must stop command substitution, and must not eat find's glob:
// the shell hands `*.txt` to find unexpanded, as it did with double quotes.
func TestSearchCommandQuotesTheTermForTheShell(t *testing.T) {
	got, err := searchFileCommand("linux", ".", "$(id)")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if want := `find '.' -type f -name '$(id)'`; got != want {
		t.Errorf("command = %q, want %q", got, want)
	}

	got, err = searchFileCommand("linux", ".", "*.txt")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if want := `find '.' -type f -name '*.txt'`; got != want {
		t.Errorf("command = %q, want %q", got, want)
	}
}
