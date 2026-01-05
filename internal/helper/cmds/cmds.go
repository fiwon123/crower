package cmds

import "github.com/spf13/cobra"

// Get name flag for cmds
func AddNameFlag(cmd *cobra.Command, name *string) {
	cmd.Flags().StringVarP(
		name,
		"name",
		"n",
		"",
		"command name")
}

// Get alias flag for cmds
func AddAllAliasFlag(cmd *cobra.Command, allAlias *[]string) {
	cmd.Flags().StringSliceVarP(
		allAlias,
		"alias",
		"a",
		[]string{},
		"define alias (--alias 'a1,a2,a3')")
}

// Get key flag for cmds
func AddKeyFlag(cmd *cobra.Command, name *string) {
	cmd.Flags().StringVarP(
		name,
		"key",
		"k",
		"",
		"command name or command alias")
}

// Get exec flag for cmds
func AddExecFlag(cmd *cobra.Command, exec *string) {
	cmd.Flags().StringVarP(
		exec,
		"exec",
		"e",
		"",
		`define the command (--exec "echo 'Hello World!'")`)
}

// Get cfg flag for cmds
func GetPersistentConfigFlag(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString("config")
}
