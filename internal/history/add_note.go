package history

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/commands"
)

func GenerateAddNote(command *commands.Data) string {
	return fmt.Sprintf("Added: %v", command.Name)
}

func GenerateAddProcessNote(command *commands.Data) string {
	return fmt.Sprintf("Added By Process: %v", command.Name)
}
