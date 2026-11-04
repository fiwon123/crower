package cmd

import (
	"testing"
	"time"
)

func TestNoticeDue(t *testing.T) {
	now := time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		mode      string
		lastCheck string
		want      bool
	}{
		{"never checked", "", "", true},
		{"daily due", "", now.Add(-25 * time.Hour).Format(time.RFC3339), true},
		{"daily not due", "", now.Add(-23 * time.Hour).Format(time.RFC3339), false},
		{"weekly due", "weekly", now.Add(-8 * 24 * time.Hour).Format(time.RFC3339), true},
		{"weekly not due", "weekly", now.Add(-6 * 24 * time.Hour).Format(time.RFC3339), false},
		{"never", "never", "", false},
		{"never despite a due stamp", "never", now.Add(-1000 * time.Hour).Format(time.RFC3339), false},
		{"unreadable stamp looks due", "", "not a time", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noticeDue(tt.mode, tt.lastCheck, now); got != tt.want {
				t.Errorf("noticeDue(%q, %q) = %v, want %v", tt.mode, tt.lastCheck, got, tt.want)
			}
		})
	}
}

func TestNoticeLine(t *testing.T) {
	if got := noticeLine("v0.8.0", "v0.8.0"); got != "" {
		t.Errorf("same version = %q, want no notice", got)
	}

	if got := noticeLine("v0.8.0", ""); got != "" {
		t.Errorf("unknown latest = %q, want no notice", got)
	}

	got := noticeLine("v0.8.0", "v0.9.0")
	if want := "crower v0.9.0 is available (you have v0.8.0) — run crower --upgrade"; got != want {
		t.Errorf("notice = %q, want %q", got, want)
	}
}
