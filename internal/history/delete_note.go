package history

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
)

func GenerateDeleteNote(command *data.Command) string {
	return fmt.Sprintf("Deleted: %v", command.Name)
}
