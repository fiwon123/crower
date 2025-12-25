package cmds

import "github.com/spf13/cobra"

func AddNameFlag(cmd *cobra.Command, name *string) {
	cmd.Flags().StringVarP(
		name,
		"name",
		"n",
		"",
		"command name")
}

func AddAllAliasFlag(cmd *cobra.Command, allAlias *[]string) {
	cmd.Flags().StringSliceVarP(
		allAlias,
		"alias",
		"a",
		[]string{},
		"define alias (--alias 'a1,a2,a3')")
}

func AddExecFlag(cmd *cobra.Command, exec *string) {
	cmd.Flags().StringVarP(
		exec,
		"exec",
		"e",
		"",
		`define the command (--exec "echo 'Hello World!'")`)
}

func GetPersistentConfigFlag(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString("config")
}
