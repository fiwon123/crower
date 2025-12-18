package handlers

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
)

func List(app *data.App) {
	fmt.Println("list all commands:")
	fmt.Println("-----------------------------")
	printMap(app.AllCommandsByName)
}

func printMap(m map[string]data.Command) {
	for k, command := range m {
		fmt.Printf("[%s] \n", k)
		fmt.Printf("  Name    = %s \n", command.Name)
		fmt.Printf("  Aliases = %s \n", command.AllAlias)
		fmt.Printf("  Exec    = %s \n", command.Exec)
		fmt.Println("")
	}
}
