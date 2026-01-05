package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func UpgradeApp(currentVersion string, app *app.Data) {
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

	err = handlers.UpgradeApp(newVersion, app)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("crower upgraded from %s to %s \n", currentVersion, newVersion)
}
