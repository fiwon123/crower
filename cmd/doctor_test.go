package cmd

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDoctorJSON(t *testing.T) {
	dir := t.TempDir()

	stdout, _, err := runWith(t, dir+"/crower.yaml", "", "doctor", "--json")
	if err != nil {
		t.Fatalf("doctor --json: %v", err)
	}

	var report map[string]any
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("not JSON: %v (%q)", err, stdout)
	}

	for _, key := range []string{"setup", "config_file", "config_folder", "config_writable", "commands", "history", "update_check", "shell", "editor", "terminal", "goos", "goarch"} {
		if _, ok := report[key]; !ok {
			t.Errorf("report is missing %q: %v", key, report)
		}
	}
}

func TestDoctorText(t *testing.T) {
	dir := t.TempDir()

	stdout, _, err := runWith(t, dir+"/crower.yaml", "", "doctor")
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}

	for _, want := range []string{"Config file:", "Config writable:", "Commands:", "Update check:", "Platform:"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("doctor output = %q, want %q", stdout, want)
		}
	}
}
