package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/command"
)

// Create a new execute note
func GenerateExecuteNote(command *command.Data) string {
	return fmt.Sprintf("Execute: %v", command.Name)
}
