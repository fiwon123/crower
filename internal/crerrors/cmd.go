package crerrors

import "fmt"

func PrintCmdHelp(cmdName string) {
	fmt.Printf("Type 'crower %s --help' for more information. \n", cmdName)
}
