package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/command"
)

func GenerateExecuteNote(command *command.Data) string {
	return fmt.Sprintf("Execute: %v", command.Name)
}
