package notes

import (
	"fmt"
)

// Create a new restore note
func GenerateRestoreNote(msg string) string {
	return fmt.Sprintf("Restore: %v", msg)
}
