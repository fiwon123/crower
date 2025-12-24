package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/command"
)

func GenerateAddNote(command *command.Data) string {
	return fmt.Sprintf("Added: %v", command.Name)
}

func GenerateAddProcessNote(command *command.Data) string {
	return fmt.Sprintf("Added By Process: %v", command.Name)
}
