package system_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
	crowererrors "github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

func TestExtractCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		ext     string
		want    string
		wantErr bool
		errText string
	}{
		{"tar", "linux", "tar", `tar -xf 'a.tar' -C 'out'`, false, ""},
		{"gz", "linux", "gz", `tar -xzf 'a.gz' -C 'out'`, false, ""},
		{"tgz", "darwin", "tgz", `tar -xzf 'a.tgz' -C 'out'`, false, ""},
		{"bz2", "linux", "bz2", `tar -xjf 'a.bz2' -C 'out'`, false, ""},
		{"xz", "linux", "xz", `tar -xJf 'a.xz' -C 'out'`, false, ""},
		{"zip on linux", "linux", "zip", `unzip 'a.zip' -d  'out'`, false, ""},
		{"zip on darwin", "darwin", "zip", `unzip 'a.zip' -d  'out'`, false, ""},
		{"zip on windows uses tar", "windows", "zip", `tar -xf a.zip -C out`, false, ""},
		{"7z", "linux", "7z", `7z x 'a.7z' -o'out'`, false, ""},
		{"rar", "windows", "rar", `7z x a.rar -oout`, false, ""},
		{"unsupported platform", "plan9", "zip", "", true, "unsupported"},
		{"unknown extension", "linux", "iso", "", true, "no extractor for .iso archives"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := system.ExtractCommand(tt.goos, tt.ext, "a."+tt.ext, "out")

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ExtractCommand(%q, %q) returned no error", tt.goos, tt.ext)
				}

				if tt.errText == "unsupported" {
					if !errors.Is(err, crowererrors.ErrUnsupportedOS) {
						t.Fatalf("error = %v, want ErrUnsupportedOS", err)
					}
				} else if !strings.Contains(err.Error(), tt.errText) {
					t.Errorf("error = %v, want it to contain %q", err, tt.errText)
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

// The archive and destination are quoted, not pasted: a name with a space or a
// quote must stay one argument.
func TestExtractCommandQuotesPaths(t *testing.T) {
	got, err := system.ExtractCommand("linux", "tar", "/tmp/my archive.tar", "/tmp/my out")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if want := `tar -xf '/tmp/my archive.tar' -C '/tmp/my out'`; got != want {
		t.Errorf("command = %q, want %q", got, want)
	}
}

// The table only proves the line is right; this proves tar runs and unpacks.
func TestArchiverExtractsATar(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test builds a tar with the host tool")
	}

	dir := t.TempDir()
	source := filepath.Join(dir, "content.txt")
	if err := os.WriteFile(source, []byte("hello"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	archive := filepath.Join(dir, "a.tar")
	if _, err := system.NewRunner().Run("tar -cf "+shell.Quote(archive)+" -C "+shell.Quote(dir)+" content.txt", system.Target{}); err != nil {
		t.Fatalf("build the tar: %v", err)
	}

	outDir := filepath.Join(dir, "out")
	if err := os.Mkdir(outDir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	out, err := system.NewArchiver().Extract("tar", archive, outDir)
	if err != nil {
		t.Fatalf("Extract: %v (output %q)", err, out)
	}

	if _, err := os.Stat(filepath.Join(outDir, "content.txt")); err != nil {
		t.Errorf("extracted file: %v", err)
	}
}
