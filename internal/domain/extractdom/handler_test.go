package extractdom_test

import (
	"archive/tar"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fiwon123/crower/internal/domain/extractdom"
	"github.com/fiwon123/crower/internal/testutil"
)

func newHandler(t *testing.T) *extractdom.Handler {
	t.Helper()

	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("InitCrowerTests: %v", err)
	}

	return extractdom.NewHandler(cfg)
}

func writeTar(t *testing.T, archivePath string, name string, content string) {
	t.Helper()

	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create %s: %v", archivePath, err)
	}

	writer := tar.NewWriter(file)

	header := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(content))}
	if err := writer.WriteHeader(header); err != nil {
		t.Fatalf("write header: %v", err)
	}

	if _, err := writer.Write([]byte(content)); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}

	// The trailer is only written on Close, and the file must be closed after
	// it: a deferred Close would hide both failures behind a confusing error
	// from Extract.
	if err := writer.Close(); err != nil {
		t.Fatalf("close tar %s: %v", archivePath, err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close %s: %v", archivePath, err)
	}
}

// prepare writes one archive into each of two sibling directories, and leaves
// the directories without the archive's file, so a misplaced extraction is
// visible.
func prepare(t *testing.T) (aDir string, bDir string, aTar string, bTar string) {
	t.Helper()

	root := t.TempDir()
	aDir = filepath.Join(root, "a")
	bDir = filepath.Join(root, "b")

	for _, dir := range []string{aDir, bDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	aTar = filepath.Join(aDir, "a.tar")
	bTar = filepath.Join(bDir, "b.tar")

	writeTar(t, aTar, "a.txt", "from a")
	writeTar(t, bTar, "b.txt", "from b")

	return aDir, bDir, aTar, bTar
}

// Without --out each archive unpacks next to itself. The flag used to be
// reassigned inside the loop, so every archive after the first went into the
// first archive's directory.
func TestExtractWithoutOutGoesNextToEachArchive(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no gate runs a Windows shell: extraction there is unverified, and QuoteFor's Windows form has known gaps")
	}

	aDir, bDir, aTar, bTar := prepare(t)

	if err := newHandler(t).Extract([]string{aTar, bTar}, ""); err != nil {
		t.Fatalf("Extract: %v", err)
	}

	if _, err := os.Stat(filepath.Join(aDir, "a.txt")); err != nil {
		t.Errorf("a.txt is not next to its archive: %v", err)
	}

	if _, err := os.Stat(filepath.Join(bDir, "b.txt")); err != nil {
		t.Errorf("b.txt is not next to its archive: %v", err)
	}

	if _, err := os.Stat(filepath.Join(aDir, "b.txt")); err == nil {
		t.Error("b.tar was extracted into a.tar's directory")
	}
}

// An explicit --out still sends every archive to that one directory.
func TestExtractWithOutUsesThatFolder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no gate runs a Windows shell: extraction there is unverified, and QuoteFor's Windows form has known gaps")
	}

	aDir, bDir, aTar, bTar := prepare(t)
	outDir := filepath.Join(t.TempDir(), "out")

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", outDir, err)
	}

	if err := newHandler(t).Extract([]string{aTar, bTar}, outDir); err != nil {
		t.Fatalf("Extract: %v", err)
	}

	for _, name := range []string{"a.txt", "b.txt"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Errorf("%s is not in the requested folder: %v", name, err)
		}
	}

	if _, err := os.Stat(filepath.Join(aDir, "a.txt")); err == nil {
		t.Error("a.txt was extracted next to its archive despite --out")
	}

	if _, err := os.Stat(filepath.Join(bDir, "b.txt")); err == nil {
		t.Error("b.txt was extracted next to its archive despite --out")
	}
}

// An archive whose name has no extension must fail: the old code skipped it
// and returned a nil error, so `crower extract something.` exited 0 having run
// nothing.
func TestExtractWithoutAnExtensionFails(t *testing.T) {
	handler := newHandler(t)
	archive := filepath.Join(t.TempDir(), "archive.")

	if err := handler.Extract([]string{archive}, ""); err == nil {
		t.Fatal("Extract accepted an archive with no extension")
	}
}

// An extension with no extractor is a failure too, and it must not reach the
// shell: the format is rejected before a command is built.
func TestExtractUnknownFormatFails(t *testing.T) {
	handler := newHandler(t)
	archive := filepath.Join(t.TempDir(), "archive.lz4")

	if err := handler.Extract([]string{archive}, ""); err == nil {
		t.Fatal("Extract accepted a format it cannot extract")
	}
}

// Every archive is attempted: one bad path must not stop the rest, and the
// failures are returned together.
func TestExtractKeepsGoingAfterAFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no gate runs a Windows shell: extraction there is unverified, and QuoteFor's Windows form has known gaps")
	}

	aDir, _, aTar, _ := prepare(t)
	bad := filepath.Join(t.TempDir(), "bad.lz4")

	if err := newHandler(t).Extract([]string{bad, aTar}, ""); err == nil {
		t.Fatal("Extract returned no error for a format it cannot extract")
	}

	if _, err := os.Stat(filepath.Join(aDir, "a.txt")); err != nil {
		t.Errorf("the good archive was not extracted after a bad one: %v", err)
	}
}

// A bare archive name has no directory component, so the default destination
// is the working directory.
func TestExtractWithoutOutInTheWorkingDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no gate runs a Windows shell: extraction there is unverified, and QuoteFor's Windows form has known gaps")
	}

	dir := t.TempDir()
	writeTar(t, filepath.Join(dir, "a.tar"), "a.txt", "from a")
	writeTar(t, filepath.Join(dir, "b.tar"), "b.txt", "from b")

	t.Chdir(dir)

	if err := newHandler(t).Extract([]string{"a.tar", "b.tar"}, ""); err != nil {
		t.Fatalf("Extract: %v", err)
	}

	for _, name := range []string{"a.txt", "b.txt"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("%s is not in the working directory: %v", name, err)
		}
	}
}
