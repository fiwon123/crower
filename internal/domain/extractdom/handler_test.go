package extractdom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/extractdom"
	"github.com/fiwon123/crower/internal/helper"
)

// An archive whose name has no extension must fail: the old code skipped it
// and returned a nil error, so `crower extract something.` exited 0 having run
// nothing.
func TestExtractWithoutAnExtensionFails(t *testing.T) {
	cfg, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("InitCrowerTests: %v", err)
	}

	handler := extractdom.NewHandler(cfg)

	archive := filepath.Join(t.TempDir(), "archive.")

	if err := handler.Extract([]string{archive}, ""); err == nil {
		t.Fatal("Extract accepted an archive with no extension")
	}
}

// An extension with no extractor is a failure too, and it must not reach the
// shell: the format is rejected before a command is built.
func TestExtractUnknownFormatFails(t *testing.T) {
	cfg, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("InitCrowerTests: %v", err)
	}

	handler := extractdom.NewHandler(cfg)

	archive := filepath.Join(t.TempDir(), "archive.lz4")

	if err := handler.Extract([]string{archive}, ""); err == nil {
		t.Fatal("Extract accepted a format it cannot extract")
	}
}
