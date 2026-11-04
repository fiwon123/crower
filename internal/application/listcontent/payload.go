// Package listcontent renders the six things `list` can report: the managed
// commands, the history, the running processes, a folder, the environment
// variables and PATH.
//
// Each mode has a text rendering for a person and a payload for `--json`. Both
// go through the Output port, and the payload types live here because they are
// the operation's result contract: the JSON keys are what a script reads.
package listcontent

import (
	"os"
	"strings"
)

// CommandRow is one row of the command table.
type CommandRow struct {
	Row         int      `json:"row"`
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases"`
	Exec        string   `json:"exec"`
	Description string   `json:"description,omitempty"`
	Params      []string `json:"params,omitempty"`
}

// CommandsResult is the result of `list` and `list --json`.
type CommandsResult struct {
	Commands []CommandRow `json:"commands"`
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

// ProcessRow is one running process.
type ProcessRow struct {
	PID  int32  `json:"pid"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// ProcessesResult is the result of `list --process --json`.
type ProcessesResult struct {
	Processes []ProcessRow `json:"processes"`
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

// FolderResult is the result of `list --folder --json`.
type FolderResult struct {
	Folder  string        `json:"folder"`
	Entries []FolderEntry `json:"entries"`
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
