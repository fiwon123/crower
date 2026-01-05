package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckNewVersion(currentVersion string, app *app.Data) {
	newVersion, err := handlers.CheckNewVersion(currentVersion, app)
	if err != nil {
		fmt.Println(err)
		return
	}

	if newVersion == currentVersion {
		fmt.Printf("Current Version: %s \n", currentVersion)
		fmt.Println("already up-to-date")
		return
	}

	fmt.Printf("Current Version: %s \n", currentVersion)
	fmt.Printf("New Version Found: %s \n", newVersion)
	fmt.Println("Check: https://github.com/fiwon123/crower/releases/latest")
}
