package tui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Executor runs a managed command for the interface.
//
// It is declared here, at the point of use, so this package does not depend on
// the use case or its adapters; the command layer supplies one.
type Executor interface {
	// Execute runs key with params and returns the output and the command's
	// name.
	Execute(key string, params []string) (string, string, error)
}

// TUI is the menu interface.
//
// It takes the catalogue and the timeline as values: the interface renders
// them, it does not own them, and it never touches the composition root.
type TUI struct {
	commands catalog.Catalog
	history  timeline.Log
	exec     Executor
	reader   *bufio.Reader
	out      io.Writer
}

// NewTUI builds the menu over the current command set and history.
func NewTUI(commands catalog.Catalog, history timeline.Log, exec Executor) *TUI {
	return newTUI(commands, history, exec, os.Stdin, os.Stdout)
}

// newTUI is the injectable constructor a test uses; NewTUI is the terminal.
func newTUI(commands catalog.Catalog, history timeline.Log, exec Executor, in io.Reader, out io.Writer) *TUI {
	return &TUI{
		commands: commands,
		history:  history,
		exec:     exec,
		reader:   bufio.NewReader(in),
		out:      out,
	}
}

func (t *TUI) Run() error {
	for {
		fmt.Fprintln(t.out, "\n╔═══════════════════════════════════════╗")
		fmt.Fprintln(t.out, "║           CROWER - TUI GUI           ║")
		fmt.Fprintln(t.out, "╚═══════════════════════════════════════╝")
		fmt.Fprintln(t.out, "\n1. List Commands")
		fmt.Fprintln(t.out, "2. Run Command")
		fmt.Fprintln(t.out, "3. Run Command with Parameters")
		fmt.Fprintln(t.out, "4. View History")
		fmt.Fprintln(t.out, "5. Export Commands")
		fmt.Fprintln(t.out, "6. Exit")
		fmt.Fprint(t.out, "\nSelect option: ")

		input, _ := t.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			t.listCommands()
		case "2":
			t.runCommand(false)
		case "3":
			t.runCommand(true)
		case "4":
			t.viewHistory()
		case "5":
			t.exportCommands()
		case "6":
			fmt.Fprintln(t.out, "Goodbye!")
			return nil
		default:
			fmt.Fprintln(t.out, "Invalid option")
		}
	}
}

func (t *TUI) listCommands() {
	fmt.Fprintln(t.out, "\n=== Commands ===")
	commands := t.commands.Commands()
	if len(commands) == 0 {
		fmt.Fprintln(t.out, "No commands found")
		return
	}

	for i, command := range commands {
		fmt.Fprintf(t.out, "%d. %s - %s\n", i+1, command.Name, command.Exec)
		if len(command.Aliases) > 0 {
			fmt.Fprintf(t.out, "   Aliases: %v\n", command.Aliases)
		}
	}
}

func (t *TUI) runCommand(withParams bool) {
	t.listCommands()
	fmt.Fprint(t.out, "\nEnter command name/number: ")
	input, _ := t.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	names := t.commands.Names()

	var cmdName string
	if idx, err := strconv.Atoi(input); err == nil && idx > 0 && idx <= len(names) {
		cmdName = names[idx-1]
	} else {
		cmdName = input
	}

	var params []string
	if withParams {
		fmt.Fprint(t.out, "Enter parameters (space separated): ")
		pinput, _ := t.reader.ReadString('\n')
		pinput = strings.TrimSpace(pinput)
		if pinput != "" {
			params = strings.Fields(pinput)
		}
	}

	fmt.Fprintln(t.out, "\n--- Executing ---")
	out, name, err := t.exec.Execute(cmdName, params)
	if name != "" {
		fmt.Fprintf(t.out, "Command: %s\n", name)
	}
	if out != "" {
		fmt.Fprintln(t.out, out)
	}
	if err != nil {
		fmt.Fprintf(t.out, "Error: %v\n", err)
	}
	fmt.Fprintln(t.out, "--- Done ---")
	fmt.Fprint(t.out, "\nPress Enter to continue...")
	t.reader.ReadString('\n')
}

func (t *TUI) viewHistory() {
	fmt.Fprintln(t.out, "\n=== History ===")
	if len(t.history.AllData) == 0 {
		fmt.Fprintln(t.out, "No history")
		return
	}
	for i, h := range t.history.AllData {
		fmt.Fprintf(t.out, "%d. [%s] %s - %s\n", i+1, h.Operation, h.CommandName, h.Timestemp)
	}
	fmt.Fprint(t.out, "\nPress Enter to continue...")
	t.reader.ReadString('\n')
}

func (t *TUI) exportCommands() {
	fmt.Fprint(t.out, "\nExport to file (default: crower-export.toml): ")
	input, _ := t.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = "crower-export.toml"
	}
	// Use sharedom logic would be better but keep it simple
	fmt.Fprintf(t.out, "Export functionality available via 'crower export' command\n")
	fmt.Fprintf(t.out, "Would export to: %s\n", input)
	fmt.Fprint(t.out, "\nPress Enter to continue...")
	t.reader.ReadString('\n')
}
