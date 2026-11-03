package extractdom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestExtractCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		ext     string
		want    string
		wantErr bool
	}{
		{"tar on linux", "linux", "tar", `tar -xf 'a.tar' -C 'out'`, false},
		{"tar on darwin", "darwin", "tar", `tar -xf 'a.tar' -C 'out'`, false},
		{"gz on linux", "linux", "gz", `tar -xzf 'a.gz' -C 'out'`, false},
		{"tgz on darwin", "darwin", "tgz", `tar -xzf 'a.tgz' -C 'out'`, false},
		{"bz2 on linux", "linux", "bz2", `tar -xjf 'a.bz2' -C 'out'`, false},
		{"xz on linux", "linux", "xz", `tar -xJf 'a.xz' -C 'out'`, false},
		{"zip on windows uses tar", "windows", "zip", `tar -xf a.zip -C out`, false},
		{"zip on linux uses unzip", "linux", "zip", `unzip 'a.zip' -d  'out'`, false},
		{"zip on darwin uses unzip", "darwin", "zip", `unzip 'a.zip' -d  'out'`, false},
		{"7z on windows", "windows", "7z", `7z x a.7z -oout`, false},
		{"7z on linux", "linux", "7z", `7z x 'a.7z' -o'out'`, false},
		{"rar on darwin", "darwin", "rar", `7z x 'a.rar' -o'out'`, false},
		{"unsupported platform", "plan9", "tar", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := "a." + tt.ext

			got, err := extractCommand(tt.goos, tt.ext, archive, "out")

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

// An extension nobody implements is a plain error, not an unsupported
// platform: the platform is fine, the archive is not.
func TestExtractCommandUnknownExtension(t *testing.T) {
	_, err := extractCommand("linux", "lz4", "a.lz4", "out")
	if err == nil {
		t.Fatal("extractCommand accepted an extension it cannot extract")
	}

	if errors.Is(err, crowererrors.ErrUnsupportedOS) {
		t.Errorf("error = %v, want an archive-format error", err)
	}
}

// Every argument is quoted for the platform that will run the line. A path
// containing the quote character must not be able to close it: a hand-written
// `'%s'` let the shell run the rest of the path.
func TestExtractCommandQuotesForThePlatform(t *testing.T) {
	tests := []struct {
		name string
		goos string
		path string
		want string
	}{
		{
			"a single quote on linux",
			"linux",
			`a'; rm -rf ~; '.tar`,
			`tar -xf 'a'\''; rm -rf ~; '\''.tar' -C 'out'`,
		},
		{
			"a single quote on darwin",
			"darwin",
			`a'; rm -rf ~; '.tar`,
			`tar -xf 'a'\''; rm -rf ~; '\''.tar' -C 'out'`,
		},
		{
			"a space on windows",
			"windows",
			`a b.tar`,
			`tar -xf "a b.tar" -C out`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractCommand(tt.goos, "tar", tt.path, "out")
			if err != nil {
				t.Fatalf("error = %v", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}
