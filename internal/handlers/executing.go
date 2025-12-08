package handlers

import (
	"fmt"
	"os/exec"

	"github.com/fiwon123/crower/internal/data"
)

func Execute(command data.Command) ([]byte, error) {

	fmt.Println(command.Exec)
	// i := strings.IndexByte(command.Exec, ' ')
	// cmdName := ""
	// args := ""
	// if i == -1 {
	// 	cmdName = command.Exec
	// } else {
	// 	cmdName = command.Exec[:i]
	// 	args = command.Exec[i+1:]
	// }

	// c := exec.Command(cmdName, args)
	c := exec.Command("sh", "-c", command.Exec)
	out, err := c.Output()

	return out, err
}
