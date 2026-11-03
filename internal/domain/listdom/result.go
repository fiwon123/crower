package listdom

import (
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/utils"
)

// The JSON results `list --json` emits.
//
// These are their own types, not tags added to the stored structs: history is
// persisted as JSON and the command map as TOML under their field names, so an
// output tag there would rename the on-disk keys and make every existing
// configuration and history file unreadable.
//
// Every collection is built with make, so an empty result encodes as [] and a
// script never has to special-case null.

// CommandRow is one row of the command table.
type CommandRow struct {
	Row     int      `json:"row"`
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
	Exec    string   `json:"exec"`
}

// CommandsResult is the result of `list` and `list --json`.
type CommandsResult struct {
	Commands []CommandRow `json:"commands"`
}

func commandsPayload(orderKeys []string, all command.MapData) CommandsResult {
	rows := make([]CommandRow, 0, len(orderKeys))

	for i, key := range orderKeys {
		single := all.Get(key)
		if single == nil {
			continue
		}

		aliases := single.AllAlias
		if aliases == nil {
			aliases = []string{}
		}

		rows = append(rows, CommandRow{
			Row:     i,
			Name:    single.Name,
			Aliases: aliases,
			Exec:    single.Exec,
		})
	}

	return CommandsResult{Commands: rows}
}

// HistoryRow is one recorded operation.
//
// Command is empty for an operation that does not touch a managed command
// (copy, move, open, list, ...).
type HistoryRow struct {
	Version   int    `json:"version"`
	File      string `json:"file"`
	Timestamp string `json:"timestamp"`
	Operation string `json:"operation"`
	Command   string `json:"command"`
	Note      string `json:"note"`
}

// HistoryResult is the result of `list --history --json`.
type HistoryResult struct {
	History []HistoryRow `json:"history"`
}

// historyPayload renders the history newest first, the order the table prints.
func historyPayload(rows []history.Content) HistoryResult {
	payload := HistoryResult{History: make([]HistoryRow, 0, len(rows))}

	for i := len(rows) - 1; i >= 0; i-- {
		row := rows[i]
		payload.History = append(payload.History, HistoryRow{
			Version:   row.Version,
			File:      row.File,
			Timestamp: row.Timestemp,
			Operation: string(row.Operation),
			Command:   row.CommandName,
			Note:      row.Note,
		})
	}

	return payload
}

// ProcessesResult is the result of `list --process --json`.
type ProcessesResult struct {
	Processes []utils.ProcessInfo `json:"processes"`
}

func processesPayload(rows []utils.ProcessInfo) ProcessesResult {
	if rows == nil {
		rows = []utils.ProcessInfo{}
	}

	return ProcessesResult{Processes: rows}
}

// processTable renders the process table `list --process` prints.
func processTable(rows []utils.ProcessInfo) string {
	const (
		titleLayout   = "%-8s %-30s %s\n"
		contentLayout = "%-8d %-30s %s\n"
	)

	var builder strings.Builder
	fmt.Fprintf(&builder, titleLayout, "PID", "Name", "Path")
	builder.WriteString("-------------------------------------------------------------\n")

	for _, row := range rows {
		fmt.Fprintf(&builder, contentLayout, row.PID, row.Name, row.Path)
	}

	return builder.String()
}

// FolderEntry is one entry of a directory listing.
type FolderEntry struct {
	Name     string `json:"name"`
	Dir      bool   `json:"dir"`
	Size     int64  `json:"size"`
	Mode     string `json:"mode"`
	Modified string `json:"modified"`
	// Error is set when the OS would not report the metadata of an entry --
	// a dangling symlink is still a directory entry.
	Error string `json:"error,omitempty"`
}

// Variable is one system variable and its OS path-list entries.
type Variable struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// VariablesResult is the result of `list --system --json`.
type VariablesResult struct {
	Variables []Variable `json:"variables"`
}

// VariableResult is the result of `list --syspath --json`.
type VariableResult struct {
	Variable Variable `json:"variable"`
}

// variablesPayload splits the raw `set`/`printenv` output the same way the text
// rendering does, so both formats describe the same variables.
func variablesPayload(raw string) VariablesResult {
	payload := VariablesResult{Variables: make([]Variable, 0)}

	for _, line := range strings.Split(raw, "\n") {
		name, value, ok := strings.Cut(strings.TrimRight(line, "\r"), "=")
		if !ok {
			continue
		}

		payload.Variables = append(payload.Variables, newVariable(name, value))
	}

	return payload
}

func variablePayload(name string, raw string) VariableResult {
	return VariableResult{Variable: newVariable(name, raw)}
}

func newVariable(name string, value string) Variable {
	// `echo` ends the value with the platform's line terminator; a carriage
	// return belongs to the line, not to the last path entry.
	value = strings.TrimRight(value, "\r\n")

	if value == "" {
		// strings.Split("", sep) is [""]: a phantom empty entry for a variable
		// that has no value at all.
		return Variable{Name: name, Values: []string{}}
	}

	return Variable{
		Name:   name,
		Values: strings.Split(value, string(os.PathListSeparator)),
	}
}
