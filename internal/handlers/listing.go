package handlers

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/data"
)

func List(app *data.App) {
	fmt.Println("------------------------------------------------")
	printMap(app.AllCommandsByName)
}

func printMap(m map[string]data.Command) {
	fmt.Printf("ID %-12s %-16s %-8s \n", "Name", "Aliases", "Exec")
	fmt.Println("------------------------------------------------")
	for _, command := range m {
		fmt.Printf("%-12s %-16v %-8s \n", command.Name, strings.Join(command.AllAlias, ","), command.Exec)
	}
}
