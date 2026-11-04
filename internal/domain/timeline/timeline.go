// Package timeline models the operation log: the ordered entries a run can
// revert to or restore from.
//
// It is a domain package: the entries, their version numbers and the rules that
// add or drop one. Rendering the table belongs to the output adapter.
package timeline

import (
	"fmt"
	"time"

	"github.com/fiwon123/crower/internal/domain/operation"
)

// Log is the operation history, oldest first.
//
// The zero value is an empty log and is usable.
type Log struct {
	AllData []Content

	// HistoryFolderPath is kept only so the history index keeps the key it
	// has always carried. Nothing reads it; dropping it would rewrite every
	// existing file's shape.
	HistoryFolderPath string
}

// Content is one recorded operation.
type Content struct {
	Version     int
	File        string
	Timestemp   string
	Operation   operation.Operation
	CommandName string
	Note        string
}

// GetBeforeLast returns the entry steps back from the end.
func (h *Log) GetBeforeLast(steps int) (*Content, error) {
	if len(h.AllData) == 0 {
		return nil, fmt.Errorf("history is empty")
	}

	if len(h.AllData)-steps < 0 {
		return nil, fmt.Errorf("step value is greater than quantity of history registry")
	}

	return &h.AllData[len(h.AllData)-steps], nil
}

// GetLast returns the newest entry, or nil for an empty log.
func (h *Log) GetLast() *Content {
	if len(h.AllData) == 0 {
		return nil
	}

	return &h.AllData[len(h.AllData)-1]
}

// GetLastOperation returns the newest entry of one kind.
func (h *Log) GetLastOperation(op operation.Operation) *Content {
	lenAllData := len(h.AllData)

	if lenAllData == 0 {
		return nil
	}

	currentIndex := lenAllData - 1
	for currentIndex >= 0 {
		if h.AllData[currentIndex].Operation == op {
			break
		}

		currentIndex -= 1
	}

	if currentIndex < 0 {
		return nil
	}

	return &h.AllData[currentIndex]
}

// Add appends an entry for one operation.
//
// commandName is the command the operation touched, or "" for an operation
// that touches files or variables rather than a managed command.
func (h *Log) Add(op operation.Operation, commandName string, note string) {
	version := 1
	if len(h.AllData) != 0 {
		lastData := h.AllData[len(h.AllData)-1]
		version = lastData.Version + 1
	}

	data := Content{
		Version:     version,
		File:        fmt.Sprintf("%05d.yaml", version),
		Timestemp:   time.Now().Format(time.RFC3339),
		Operation:   op,
		CommandName: commandName,
		Note:        note,
	}

	h.AllData = append(h.AllData, data)
}

// RemoveLast drops the newest entry.
func (h *Log) RemoveLast() {
	if len(h.AllData) == 0 {
		return
	}

	h.AllData = h.AllData[:len(h.AllData)-1]
}

// RemoveContent drops one entry.
//
// It reports whether the entry was found, and a miss leaves the log untouched.
// It used to assign an accumulator that only a match filled, so a stale entry
// replaced the whole log with nil and still returned true (#218).
func (h *Log) RemoveContent(content Content) bool {
	for i, c := range h.AllData {
		if c.Version != content.Version {
			continue
		}

		h.AllData = append(h.AllData[:i], h.AllData[i+1:]...)
		return true
	}

	return false
}

// GetIndexFromLastTo is the index steps back from the end of the log.
func (h *Log) GetIndexFromLastTo(steps int) int {
	return len(h.AllData) - steps
}
