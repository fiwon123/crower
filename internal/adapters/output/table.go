package output

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
)

// CommandsTable renders the managed command table, header included.
//
// It returns text rather than printing it because the same table is also a
// prompt: delete, execute, update and reset show it before asking which row to
// use, and there it is a diagnostic. A prompt that wrote to the result stream
// put a table on stdout during a failed run, where a caller piping the results
// read it as data.
func CommandsTable(commands catalog.Catalog) string {
	var builder strings.Builder

	builder.WriteString("------------------------------------------------\n")
	fmt.Fprintf(&builder, "%-3s %-12s %-16s %s\n", "Row", "Name", "Aliases", "Exec")
	builder.WriteString("------------------------------------------------\n")

	for i, command := range commands.Commands() {
		fmt.Fprintf(&builder, "%-3d %-12s %-16s %s\n", i, command.Name, strings.Join(command.Aliases, ","), command.Exec)
	}

	return builder.String()
}
