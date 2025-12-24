package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/command"
)

func GenerateDeleteNote(command *command.Data) string {
	return fmt.Sprintf("Deleted: %v", command.Name)
}
