package extractdom_test

import (
	"archive/tar"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fiwon123/crower/internal/domain/extractdom"
	"github.com/fiwon123/crower/internal/helper"
)

func newHandler(t *testing.T) *extractdom.Handler {
	t.Helper()

	cfg, _, err := helper.InitCrowerTests()
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
	defer file.Close()

	writer := tar.NewWriter(file)
	defer writer.Close()

	header := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(content))}
	if err := writer.WriteHeader(header); err != nil {
		t.Fatalf("write header: %v", err)
	}

	if _, err := writer.Write([]byte(content)); err != nil {
		t.Fatalf("write %s: %v", name, err)
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
		t.Skip("the extractor is the platform's tar")
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
		t.Skip("the extractor is the platform's tar")
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
