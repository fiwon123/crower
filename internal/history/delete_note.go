package history

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/commands"
)

func GenerateDeleteNote(command *commands.Data) string {
	return fmt.Sprintf("Deleted: %v", command.Name)
}
