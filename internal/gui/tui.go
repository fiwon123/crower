package gui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
)

type TUI struct {
	app     *app.Config
	exec    *executedom.Handler
	reader  *bufio.Reader
}

func NewTUI(app *app.Config) *TUI {
	return &TUI{
		app:    app,
		exec:   executedom.NewHandler(app),
		reader: bufio.NewReader(os.Stdin),
	}
}

func (t *TUI) Run() error {
	for {
		fmt.Println("\n╔═══════════════════════════════════════╗")
		fmt.Println("║           CROWER - TUI GUI           ║")
		fmt.Println("╚═══════════════════════════════════════╝")
		fmt.Println("\n1. List Commands")
		fmt.Println("2. Run Command")
		fmt.Println("3. Run Command with Parameters")
		fmt.Println("4. View History")
		fmt.Println("5. Export Commands")
		fmt.Println("6. Exit")
		fmt.Print("\nSelect option: ")

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
			fmt.Println("Goodbye!")
			return nil
		default:
			fmt.Println("Invalid option")
		}
	}
}

func (t *TUI) listCommands() {
	fmt.Println("\n=== Commands ===")
	keys := t.app.OrderKeys
	if len(keys) == 0 {
		fmt.Println("No commands found")
		return
	}
	for i, k := range keys {
		cmd := t.app.AllCommandsByName.Get(k)
		if cmd != nil {
			fmt.Printf("%d. %s - %s\n", i+1, cmd.Name, cmd.Exec)
			if len(cmd.AllAlias) > 0 {
				fmt.Printf("   Aliases: %v\n", cmd.AllAlias)
			}
		}
	}
}

func (t *TUI) runCommand(withParams bool) {
	t.listCommands()
	fmt.Print("\nEnter command name/number: ")
	input, _ := t.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var cmdName string
	if idx, err := strconv.Atoi(input); err == nil && idx > 0 && idx <= len(t.app.OrderKeys) {
		cmdName = t.app.OrderKeys[idx-1]
	} else {
		cmdName = input
	}

	var params []string
	if withParams {
		fmt.Print("Enter parameters (space separated): ")
		pinput, _ := t.reader.ReadString('\n')
		pinput = strings.TrimSpace(pinput)
		if pinput != "" {
			params = strings.Fields(pinput)
		}
	}

	fmt.Println("\n--- Executing ---")
	out, cmd, err := t.exec.Execute(cmdName, params)
	if cmd != nil {
		fmt.Printf("Command: %s\n", cmd.Name)
	}
	if out != "" {
		fmt.Println(out)
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("--- Done ---")
	fmt.Print("\nPress Enter to continue...")
	t.reader.ReadString('\n')
}

func (t *TUI) viewHistory() {
	fmt.Println("\n=== History ===")
	if len(t.app.History.AllData) == 0 {
		fmt.Println("No history")
		return
	}
	for i, h := range t.app.History.AllData {
		fmt.Printf("%d. [%s] %s - %s\n", i+1, h.Operation, h.CommandName, h.Timestemp)
	}
	fmt.Print("\nPress Enter to continue...")
	t.reader.ReadString('\n')
}

func (t *TUI) exportCommands() {
	fmt.Print("\nExport to file (default: crower-export.toml): ")
	input, _ := t.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = "crower-export.toml"
	}
	// Use sharedom logic would be better but keep it simple
	fmt.Printf("Export functionality available via 'crower export' command\n")
	fmt.Printf("Would export to: %s\n", input)
	fmt.Print("\nPress Enter to continue...")
	t.reader.ReadString('\n')
}
