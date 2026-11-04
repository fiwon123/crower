package storage_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

func TestSettingsRoundTrip(t *testing.T) {
	settings := storage.NewSettings(t.TempDir())

	if got := settings.Load(); got.UpdateCheck != "" || got.LastCheck != "" {
		t.Errorf("empty settings = %+v, want the zero value", got)
	}

	want := storage.SettingsData{UpdateCheck: "weekly", LastCheck: "2026-11-04T00:00:00Z", Latest: "v0.9.0"}
	if err := settings.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if got := settings.Load(); got != want {
		t.Errorf("Load = %+v, want %+v", got, want)
	}
}
