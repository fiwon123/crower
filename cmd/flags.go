package cmd

import "github.com/spf13/cobra"

// The flag helpers live next to the commands they configure. They used to sit in
// internal/helper, which made the helper package the only domain-free package
// importing cobra.

// AddNameFlag adds the --name flag used by create and update.
func AddNameFlag(cmd *cobra.Command, name *string) {
	cmd.Flags().StringVarP(name, "name", "n", "", "command name")
}

// AddAllAliasFlag adds the --alias flag, accepting a comma-separated list.
func AddAllAliasFlag(cmd *cobra.Command, allAlias *[]string) {
	cmd.Flags().StringSliceVarP(allAlias, "alias", "a", []string{}, "define alias (--alias 'a1,a2,a3')")
}

// AddExecFlag adds the --exec flag.
func AddExecFlag(cmd *cobra.Command, exec *string) {
	cmd.Flags().StringVarP(exec, "exec", "e", "", `define the command (--exec "echo 'Hello World!'")`)
}
