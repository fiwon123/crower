package cmd

import (
	"strings"

	"github.com/fiwon123/crower/internal/application/updatecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"

	"github.com/spf13/cobra"
)

var name string
var exec string

var create bool
var execute bool

// Cmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [command]",
	Short: "Change a command's name, aliases, exec or execution metadata.",
	Long: `Change a command's name, aliases, exec or execution metadata.

	crower update greet --exec "echo hi"
	crower update greet --name hello --alias hi

The command is found by name or by alias. The last touched command can be
updated instead:
	crower update --last
	crower update --create
	crower update --execute

With no field given, crower asks for each one. --description, --params, --cwd,
--shell and --timeout change the optional metadata, and an empty value clears it.`,
	Example: `	crower update greet --exec "echo hi"
	crower update hi --name hello
	crower update --last
	crower update --create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		session := newSessionPorts(cfg)

		request := updatecommand.Request{
			Changes: catalog.Changes{
				Name:    name,
				Aliases: allAlias,
				Exec:    exec,
			},
		}

		// A flag that was not given leaves the field alone; giving it an empty
		// value is how a user removes an icon or a group.
		if cmd.Flags().Changed("icon") {
			request.Changes.Icon = strings.TrimSpace(icon)
			request.Changes.ClearIcon = request.Changes.Icon == ""
		}

		if cmd.Flags().Changed("group") {
			request.Changes.Group = strings.TrimSpace(group)
			request.Changes.ClearGroup = request.Changes.Group == ""
		}

		// The metadata fields are pointers in catalog.Changes: a flag that was
		// given supplies the value, and an empty one clears it.
		if cmd.Flags().Changed("description") {
			value := strings.TrimSpace(description)
			request.Changes.Description = &value
		}

		if cmd.Flags().Changed("params") {
			list := flagList(params)
			request.Changes.Params = &list
		}

		if cmd.Flags().Changed("cwd") {
			value := strings.TrimSpace(cwd)
			request.Changes.Cwd = &value
		}

		if cmd.Flags().Changed("shell") {
			value := strings.TrimSpace(shell)
			request.Changes.Shell = &value
		}

		if cmd.Flags().Changed("timeout") {
			value := strings.TrimSpace(timeout)
			request.Changes.Timeout = &value
		}

		if len(args) != 0 {
			request.Key = args[0]
		}

		switch {
		case last:
			request.Last = operation.Update
		case create:
			request.Last = operation.Create
		case execute:
			request.Last = operation.Execute
		}

		return updatecommand.New(session, session, session, session).Run(request)
	},
}

func init() {
	AddNameFlag(updateCmd, &name)
	AddAllAliasFlag(updateCmd, &allAlias)
	AddExecFlag(updateCmd, &exec)
	AddIconFlag(updateCmd, &icon)
	AddGroupFlag(updateCmd, &group)
	AddMetadataFlags(updateCmd, &description, &params, &cwd, &shell, &timeout)

	updateCmd.Flags().BoolVarP(&last, "last", "l", false, "update recent updated command")
	updateCmd.Flags().BoolVarP(&create, "create", "c", false, "update recent created command")
	updateCmd.Flags().BoolVarP(&execute, "execute", "x", false, "update recent updated command")
}
