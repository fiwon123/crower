package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/command"
)

// Create a new add note
func GenerateAddNote(command *command.Data) string {
	return fmt.Sprintf("Added: %v", command.Name)
}

// Create a new add process note
func GenerateAddProcessNote(command *command.Data) string {
	return fmt.Sprintf("Added By Process: %v", command.Name)
}
