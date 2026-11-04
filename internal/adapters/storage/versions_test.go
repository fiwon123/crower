package storage_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

func TestVersionsRecordAndPrevious(t *testing.T) {
	versions := storage.NewVersions(t.TempDir())

	if _, ok := versions.Previous("v1.0.0"); ok {
		t.Fatal("an empty history reported a previous version")
	}

	for _, version := range []string{"v1.0.0", "v1.1.0"} {
		if err := versions.RecordInstalled(version); err != nil {
			t.Fatalf("RecordInstalled(%s): %v", version, err)
		}
	}

	// Newest first, so the previous version of the running one is v1.0.0.
	got, ok := versions.Previous("v1.1.0")
	if !ok || got != "v1.0.0" {
		t.Errorf("Previous(v1.1.0) = %q, %v; want v1.0.0", got, ok)
	}
}

func TestVersionsRecordMovesToTheFrontAndDeduplicates(t *testing.T) {
	versions := storage.NewVersions(t.TempDir())

	for _, version := range []string{"v1.0.0", "v1.1.0", "v1.0.0"} {
		if err := versions.RecordInstalled(version); err != nil {
			t.Fatalf("RecordInstalled(%s): %v", version, err)
		}
	}

	if got, _ := versions.Previous("v1.1.0"); got != "v1.0.0" {
		t.Errorf("Previous = %q, want v1.0.0", got)
	}

	if got, _ := versions.Previous("v2.0.0"); got != "v1.0.0" {
		t.Errorf("Previous from an unknown running version = %q, want the newest", got)
	}
}

// The list is bounded: it is a convenience, not an audit log.
func TestVersionsAreBounded(t *testing.T) {
	versions := storage.NewVersions(t.TempDir())

	for i := 0; i < 25; i++ {
		if err := versions.RecordInstalled("v1.0." + string(rune('a'+i))); err != nil {
			t.Fatalf("RecordInstalled: %v", err)
		}
	}

	// The oldest entries fell off: Previous from the newest returns the second
	// newest, not the first ever recorded.
	got, _ := versions.Previous("v1.0.y")
	if got == "v1.0.a" {
		t.Errorf("Previous = %q, want a bounded history", got)
	}

}
