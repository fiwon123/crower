package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/pkg/utils"
)

// Settings holds the small preferences that are not commands: today, whether
// crower may look for a newer release, and when it last looked.
type Settings struct {
	path string
}

// SettingsData is the file's shape.
type SettingsData struct {
	// UpdateCheck is "daily" (or empty, the default), "weekly" or "never".
	UpdateCheck string `json:"update_check,omitempty"`

	// LastCheck is when the notice last looked, RFC 3339.
	LastCheck string `json:"last_check,omitempty"`

	// Latest is the tag the last successful look found.
	Latest string `json:"latest,omitempty"`
}

// NewSettings addresses the settings file beside the configuration.
func NewSettings(folder string) *Settings {
	return &Settings{path: filepath.Join(folder, "settings.json")}
}

// Load reads the settings; a missing or unreadable file is the default, because
// a preference nobody set must not break a run.
func (s *Settings) Load() SettingsData {
	var data SettingsData

	body, err := os.ReadFile(s.path)
	if err != nil {
		return SettingsData{}
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return SettingsData{}
	}

	return data
}

// Save writes the settings atomically.
func (s *Settings) Save(data SettingsData) error {
	return utils.WriteJson(data, s.path)
}
