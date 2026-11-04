package timeline

import (
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/operation"
)

func TestAddNumbersAndNamesTheSnapshots(t *testing.T) {
	var log Log

	log.Add(operation.Create, "greet", "one")
	log.Add(operation.Copy, "", "two")

	if len(log.AllData) != 2 {
		t.Fatalf("log holds %d entries, want 2", len(log.AllData))
	}

	if log.AllData[0].Version != 1 || log.AllData[0].File != "00001.yaml" {
		t.Errorf("first = %+v, want version 1 and 00001.yaml", log.AllData[0])
	}

	if log.AllData[1].Version != 2 || log.AllData[1].File != "00002.yaml" {
		t.Errorf("second = %+v, want version 2 and 00002.yaml", log.AllData[1])
	}
}

// The entry records which command the operation touched. Without it, the
// `--last` family resolved an empty name and could never act.
func TestAddRecordsTheCommandName(t *testing.T) {
	var log Log
	log.Add(operation.Execute, "greet", "main_operation=execute;command_name=greet")

	last := log.GetLast()
	if last == nil {
		t.Fatal("Add did not append an entry")
	}

	if last.CommandName != "greet" {
		t.Errorf("CommandName = %q, want greet", last.CommandName)
	}
}

// One entry is dropped and the rest keep their order.
func TestRemoveContentDropsOneEntry(t *testing.T) {
	var log Log
	log.Add(operation.Create, "first", "one")
	log.Add(operation.Update, "second", "two")
	log.Add(operation.Delete, "third", "three")

	if ok := log.RemoveContent(log.AllData[1]); !ok {
		t.Fatal("RemoveContent of a present entry returned false")
	}

	remaining := []string{}
	for _, entry := range log.AllData {
		remaining = append(remaining, entry.CommandName)
	}

	if want := []string{"first", "third"}; !slices.Equal(remaining, want) {
		t.Errorf("history = %v, want %v", remaining, want)
	}
}

// A stale entry that is not in the log is not a reason to discard the log.
func TestRemoveContentLeavesTheLogOnAMiss(t *testing.T) {
	var log Log
	log.Add(operation.Create, "first", "one")
	log.Add(operation.Update, "second", "two")

	before := slices.Clone(log.AllData)

	if ok := log.RemoveContent(Content{Version: 99}); ok {
		t.Fatal("RemoveContent of a missing entry returned true")
	}

	if !slices.Equal(log.AllData, before) {
		t.Errorf("history changed on a miss: %v", log.AllData)
	}
}

func TestGetBeforeLast(t *testing.T) {
	var log Log

	if _, err := log.GetBeforeLast(1); err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("error = %v, want the empty message", err)
	}

	log.Add(operation.Create, "a", "one")
	log.Add(operation.Update, "b", "two")

	entry, err := log.GetBeforeLast(2)
	if err != nil {
		t.Fatalf("GetBeforeLast: %v", err)
	}

	if entry.CommandName != "a" {
		t.Errorf("entry = %+v, want the first", entry)
	}

	if _, err := log.GetBeforeLast(3); err == nil || !strings.Contains(err.Error(), "greater than") {
		t.Fatalf("error = %v, want the range message", err)
	}
}

func TestGetLastOperation(t *testing.T) {
	var log Log
	log.Add(operation.Create, "a", "one")
	log.Add(operation.Execute, "b", "two")
	log.Add(operation.Create, "c", "three")

	if got := log.GetLastOperation(operation.Create); got == nil || got.CommandName != "c" {
		t.Errorf("last create = %+v, want c", got)
	}

	if got := log.GetLastOperation(operation.Delete); got != nil {
		t.Errorf("last delete = %+v, want nil", got)
	}
}

func TestRemoveLast(t *testing.T) {
	var log Log
	log.RemoveLast() // empty is a no-op

	log.Add(operation.Create, "a", "one")
	log.Add(operation.Update, "b", "two")
	log.RemoveLast()

	if len(log.AllData) != 1 || log.AllData[0].CommandName != "a" {
		t.Errorf("log = %+v, want only the first entry", log.AllData)
	}
}
