package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/pkg/utils"
)

// Versions records the releases the running installation has installed.
//
// It exists for one question: `--rollback` has to know what was installed
// before this version. The list is newest first and bounded, because it is a
// convenience, not an audit log.
type Versions struct {
	path string
}

// maxRecordedVersions bounds the history the file keeps.
const maxRecordedVersions = 10

type versionsFile struct {
	Installed []string `json:"installed"`
}

// NewVersions addresses the version list beside the configuration file.
func NewVersions(folder string) *Versions {
	return &Versions{path: filepath.Join(folder, "versions.json")}
}

// RecordInstalled puts version at the front, dropping older duplicates.
func (v *Versions) RecordInstalled(version string) error {
	if version == "" {
		return nil
	}

	list := append([]string{version}, v.load().Installed...)

	seen := map[string]bool{}
	kept := make([]string, 0, len(list))

	for _, item := range list {
		if item == "" || seen[item] {
			continue
		}

		seen[item] = true
		kept = append(kept, item)

		if len(kept) == maxRecordedVersions {
			break
		}
	}

	return utils.WriteJson(versionsFile{Installed: kept}, v.path)
}

// Previous returns the most recently installed version other than current.
func (v *Versions) Previous(current string) (string, bool) {
	for _, version := range v.load().Installed {
		if version != "" && version != current {
			return version, true
		}
	}

	return "", false
}

// load reads the list; a missing or unreadable file is an empty list, because
// not knowing the history must not break an update.
func (v *Versions) load() versionsFile {
	var file versionsFile

	body, err := os.ReadFile(v.path)
	if err != nil {
		return versionsFile{}
	}

	if err := json.Unmarshal(body, &file); err != nil {
		return versionsFile{}
	}

	return file
}
