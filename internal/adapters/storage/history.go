package storage

import (
	"fmt"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"os"
	"path/filepath"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/storage/command"

	"github.com/fiwon123/crower/pkg/utils"
)

const (
	// indexName is the history index file. The name is historical and is part
	// of the format: renaming it would orphan every existing installation.
	indexName = "historyhelper.json"

	// snapshotFolder holds one file per operation.
	snapshotFolder = "history"
)

// HistoryFiles is the history index and the snapshot folder beside it.
type HistoryFiles struct {
	index  string
	folder string
}

// NewHistoryFiles addresses the history files of the installation stored in
// folder -- the same folder the command map lives in.
func NewHistoryFiles(folder string) *HistoryFiles {
	return &HistoryFiles{
		index:  filepath.Join(folder, indexName),
		folder: filepath.Join(folder, snapshotFolder),
	}
}

// Ensure creates the index file and the snapshot folder when they are missing.
func (h *HistoryFiles) Ensure() error {
	if err := utils.CreateFileIfNotExists(h.index); err != nil {
		return fmt.Errorf("prepare history file %s: %w", h.index, err)
	}

	if err := utils.CreateFolderIfNotExists(h.folder); err != nil {
		return fmt.Errorf("prepare history folder %s: %w", h.folder, err)
	}

	return nil
}

// Load reads the history index. An empty file is a fresh installation, not a
// broken one.
func (h *HistoryFiles) Load() (timeline.Log, error) {
	var data timeline.Log

	if err := utils.ReadJson(h.index, &data); err != nil {
		return timeline.Log{}, fmt.Errorf("read history %s: %w", h.index, err)
	}

	return data, nil
}

// SaveIndex writes the history index without taking a snapshot.
func (h *HistoryFiles) SaveIndex(data timeline.Log) error {
	return utils.WriteJson(data, h.index)
}

// WriteSnapshot stores the command map under one history entry's file name.
func (h *HistoryFiles) WriteSnapshot(file string, commands command.MapData) error {
	path, err := h.snapshotPath(file)
	if err != nil {
		return err
	}

	if err := utils.CreateFileIfNotExists(path); err != nil {
		return err
	}

	return utils.WriteToml(commands, path)
}

// ReadSnapshot reads the command map an entry points at.
func (h *HistoryFiles) ReadSnapshot(file string) (command.MapData, error) {
	path, err := h.snapshotPath(file)
	if err != nil {
		return nil, err
	}

	commands := command.NewMapData()

	if err := utils.ReadToml(path, &commands); err != nil {
		return nil, err
	}

	return commands, nil
}

// HasSnapshot reports whether an entry's file is present.
func (h *HistoryFiles) HasSnapshot(file string) bool {
	path, err := h.snapshotPath(file)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

// DeleteSnapshot removes the file one history entry points at.
func (h *HistoryFiles) DeleteSnapshot(file string) error {
	path, err := h.snapshotPath(file)
	if err != nil {
		return err
	}

	return utils.DeleteFile(path)
}

// Folder is the directory the snapshots live in.
func (h *HistoryFiles) Folder() string {
	return h.folder
}

// snapshotPath resolves an entry's file name inside the snapshot folder.
//
// file comes from a history entry, not from a caller: a history file is JSON in
// a user's home directory and may have been edited by hand, so the name is
// checked here rather than trusted. Both separators are rejected on every
// platform, not just the host's: a backslash is a legal character in a Linux
// file name but a separator on Windows, and a snapshot is always written by
// this program as `<version>.yaml`.
func (h *HistoryFiles) snapshotPath(file string) (string, error) {
	if file == "" || filepath.Base(file) != file || strings.ContainsAny(file, `/\`) {
		return "", fmt.Errorf("invalid history file name %q", file)
	}

	return filepath.Join(h.folder, file), nil
}
