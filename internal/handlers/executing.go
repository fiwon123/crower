package handlers

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fiwon123/crower/internal/data"
)

func Execute(command data.Command) ([]byte, error) {

	var c *exec.Cmd
	fmt.Println(command.Exec)
	switch runtime.GOOS {
	case "windows":
		c = exec.Command("cmd", "/c", command.Exec)
	case "linux":
		c = exec.Command("sh", "-c", command.Exec)
	}

	out, err := c.Output()
	return out, err
}
