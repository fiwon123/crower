package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

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

// AddIconFlag adds the --icon flag used by create and update.
func AddIconFlag(cmd *cobra.Command, icon *string) {
	cmd.Flags().StringVar(icon, "icon", "", "short marker shown beside the command (an emoji)")
}

// AddGroupFlag adds the --group flag used by create and update.
func AddGroupFlag(cmd *cobra.Command, group *string) {
	cmd.Flags().StringVar(group, "group", "", "name of the group the command belongs to")
}

// AddMetadataFlags adds the optional execution metadata create and update share.
// Every flag is independent: a flag that is not given leaves the field alone.
func AddMetadataFlags(cmd *cobra.Command, description *string, params *[]string, cwd, shell, timeout *string) {
	cmd.Flags().StringVar(description, "description", "", "what the command is for")
	cmd.Flags().StringSliceVar(params, "params", []string{}, "positional parameters, comma-separated or repeated (--params 'env,region')")
	cmd.Flags().StringVar(cwd, "cwd", "", "working directory the command runs in")
	cmd.Flags().StringVar(shell, "shell", "", "shell that runs the command (sh, pwsh or cmd)")
	cmd.Flags().StringVar(timeout, "timeout", "", "how long the command may run, as a Go duration (120s)")
}

// flagList trims a slice-valued flag and drops the empty entries a trailing
// comma leaves, so `--params ""` means "clear it".
func flagList(values []string) []string {
	list := make([]string, 0, len(values))

	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			list = append(list, trimmed)
		}
	}

	return list
}

// AddExecFlag adds the --exec flag.
func AddExecFlag(cmd *cobra.Command, exec *string) {
	cmd.Flags().StringVarP(exec, "exec", "e", "", `define the command (--exec "echo 'Hello World!'")`)
}
