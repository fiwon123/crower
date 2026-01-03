package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/command"
)

// Create a new delete note
func GenerateDeleteNote(command *command.Data) string {
	return fmt.Sprintf("Deleted: %v", command.Name)
}
