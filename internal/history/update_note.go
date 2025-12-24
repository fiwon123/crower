package history

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fiwon123/crower/internal/data/commands"
)

func GenerateUpdateNote(oldCommand *commands.Data, newCommand *commands.Data) string {

	output := strings.Builder{}
	output.WriteString("Updated:")
	if oldCommand.Name != newCommand.Name {
		changes := fmt.Sprintf(" name from %s to %s,", oldCommand.Name, newCommand.Name)
		output.WriteString(changes)
	}

	if hasDiffAlieses(oldCommand.AllAlias, newCommand.AllAlias) {
		changes := fmt.Sprintf(" aliases from %v to %v,", oldCommand.AllAlias, newCommand.AllAlias)
		output.WriteString(changes)
	}

	if oldCommand.Exec != newCommand.Exec {
		changes := fmt.Sprintf(" exec from \"%s\" to \"%s\"", oldCommand.Exec, newCommand.Exec)
		output.WriteString(changes)
	}

	return output.String()
}

func hasDiffAlieses(old []string, new []string) bool {
	lenOld := len(old)
	if lenOld != len(new) {
		return true
	}

	sort.Strings(old)
	sort.Strings(new)

	for i := range lenOld {
		if old[i] != new[i] {
			return true
		}
	}

	return false
}
