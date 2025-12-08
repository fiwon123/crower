package configuration

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func Open(cfgFilePath string) {
	handlers.Execute(*data.NewCommand(
		"open",
		[]string{},
		fmt.Sprintf("xdg-open %s", cfgFilePath)))
}
