//go:build linux

package handlers

import (
	"fmt"
	"os"

	"github.com/fiwon123/crower/internal/data/app"
)

func CreateSystemPathVariable(value string, app *app.Data) (string, error) {
	line := fmt.Sprintf(`export PATH="$PATH:%s"`, value)

	err := os.WriteFile(
		os.Getenv("HOME")+"/.profile",
		[]byte(line+"\n"),
		0644,
	)

	return "Added to PATH", err
}
