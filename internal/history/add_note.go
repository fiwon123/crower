package history

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
)

func GenerateAddNote(command *data.Command) string {
	return fmt.Sprintf("Added: %v", command.Name)
}

func GenerateAddProcessNote(command *data.Command) string {
	return fmt.Sprintf("Added By Process: %v", command.Name)
}
