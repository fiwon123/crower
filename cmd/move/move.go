package move

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd represents the move command
var Cmd = &cobra.Command{
	Use:   "move",
	Short: "move file or folder to other location",
	Long:  `move file or folder to other location`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("move called")
	},
}

func init() {

}
