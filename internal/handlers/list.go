package handlers

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/data"
)

func List(app *data.App) {
	fmt.Println("------------------------------------------------")
	print(app.OrderKeys, app.AllCommandsByName)
}

func print(orderKeys []string, allCommands data.CommandsMap) {
	fmt.Printf("%-3s %-12s %-16s %-8s \n", "Row", "Name", "Aliases", "Exec")
	fmt.Println("------------------------------------------------")

	for i, key := range orderKeys {
		command := allCommands.Get(key)
		fmt.Printf("%-3d %-12s %-16v %-8s \n", i, command.Name, strings.Join(command.AllAlias, ","), command.Exec)
	}
}
