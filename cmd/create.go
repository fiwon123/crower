package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/createcommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/fiwon123/crower/pkg/utils"
	"github.com/spf13/cobra"
)

var allAlias []string

var scriptFlag bool
var processName string

// Cmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [name] [exec]",
	Short: "Create a command, a file or a folder.",
	Long: `Create a command, a file or a folder.

A command has a name, optional aliases, and the line it runs:
	crower create greet "echo hello"
	crower create greet "echo hello" --alias hi

It can come from a running process, or be a script:
	crower create -p server nginx
	crower create --script deploy

A file or a folder is created at the path you give:
	crower create --file notes.txt
	crower create --folder reports

With no arguments crower asks for each field. --execute makes a command from
the last command you ran.`,
	Example: `	crower create
	crower create greet "echo hello"
	crower create greet "echo hello" --alias hi
	crower create --file notes.txt
	crower create --folder reports
	crower create -p server nginx
	crower create --script deploy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		request := createcommand.Request{
			Mode:    createcommand.ModeCommand,
			Aliases: allAlias,
			Args:    args,
		}

		// The flags decide which job this run is -- and their precedence is the
		// one create has always had -- while the use case owns what each job
		// does. A positional path keeps the convenience of creating a file or
		// folder without a flag.
		switch {
		case processName != "":
			request.Mode, request.Name = createcommand.ModeProcess, processName
		case fileFlag:
			request.Mode = createcommand.ModeFile
		case folderFlag:
			request.Mode = createcommand.ModeFolder
		case sysPathFlag:
			request.Mode = createcommand.ModeSystemPath
		case systemFlag:
			request.Mode = createcommand.ModeSystemVariable
		case executeFlag:
			request.Mode, request.Last = createcommand.ModeLast, operation.Execute
		case scriptFlag:
			request.Mode = createcommand.ModeScript
		case len(args) == 0:
			request.Mode = createcommand.ModeCommand
		case utils.IsValidFilePath(args[0]):
			request.Mode = createcommand.ModeFile
		case utils.IsValidFolderPath(args[0]):
			request.Mode = createcommand.ModeFolder
		default:
			request.Mode = createcommand.ModeCommand
		}

		session := newCreateSession(cfg)

		return createcommand.New(
			session, session, session, session, session, session, session, session, session,
		).Run(request)
	},
}

// createSession adapts the configuration to the ports the create use case
// declares.
//
// It embeds legacySession for the ports create shares with every other use
// case -- the catalog, the diagnostics stream, the confirmation -- and adds the
// rest: the prompt texts create has always used, the platform adapters, the last
// operation read from the history, and the notes.
type createSession struct {
	legacySession

	creator     *system.Creator
	scripts     *system.Scripts
	processes   *system.Processes
	environment *system.Environment
	opener      *system.Opener
}

// newCreateSession wires the configuration, the renderer and the adapters
// create needs.
func newCreateSession(cfg *app.Config) createSession {
	return createSession{
		legacySession: newLegacySession(cfg),
		creator:       system.NewCreator(),
		scripts:       system.NewScripts(filepath.Join(cfg.ConfigFolder(), "scripts")),
		processes:     system.NewProcesses(),
		environment:   system.NewEnvironment(),
		opener:        system.NewOpener(),
	}
}

// Save replaces the command set and writes the command file.
//
// The history entry is recorded separately: `create --script` has a step
// between the two.
func (s createSession) Save(commands catalog.Catalog) error {
	s.cfg.ReplaceCatalog(commands)

	return s.cfg.SaveCommands()
}

// LastOperation describes the last operation of a kind from the history note.
//
// A note written before command_steps existed carries none; the command may
// still be in the catalog, where its steps are visible.
func (s createSession) LastOperation(op operation.Operation) (createcommand.LastRun, bool) {
	content := s.cfg.GetLast(op)
	if content == nil {
		return createcommand.LastRun{}, false
	}

	steps, hasSteps := notes.CommandSteps(content.Note)

	if !hasSteps && content.CommandName != "" {
		if command := s.cfg.AllCommandsByName.Get(content.CommandName); command != nil && len(command.Steps) > 0 {
			steps, hasSteps = command.Steps, true
		}
	}

	return createcommand.LastRun{
		CommandName: content.CommandName,
		Exec:        notes.CommandExec(content.Note),
		Steps:       steps,
		HasSteps:    hasSteps,
	}, true
}

// CreateName asks for the new command's name. The text is the one create has
// always shown; the update form is a different question.
func (s createSession) CreateName() (string, error) {
	value, err := helper.GetUserInput("New Name ", helper.IsValidInput, s.cfg)
	if err != nil {
		return "", err
	}

	name, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("create name: expected a string, got %T", value)
	}

	return name, nil
}

// CreateExec asks for the line the new command runs.
func (s createSession) CreateExec() (string, error) {
	value, err := helper.GetUserInput("New Exec ", helper.IsValidInput, s.cfg)
	if err != nil {
		return "", err
	}

	exec, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("create exec: expected a string, got %T", value)
	}

	return exec, nil
}

// File creates one file with the host's own tool.
func (s createSession) File(path string) error {
	return s.creator.File(path)
}

// Folder creates one folder with the host's own tool.
func (s createSession) Folder(path string) error {
	return s.creator.Folder(path)
}

// Create writes the script file for name.
func (s createSession) Create(name string) (string, string, error) {
	return s.scripts.Create(name)
}

// ExecLine resolves a process name or pid to the line that starts it again.
func (s createSession) ExecLine(process string) (string, error) {
	return s.processes.ExecLine(process)
}

// SetVariable adds one variable to the file the shell reads at start.
func (s createSession) SetVariable(name, value string) (string, error) {
	return s.environment.SetVariable(name, value)
}

// AddPath adds one directory to PATH.
func (s createSession) AddPath(path string) (string, error) {
	return s.environment.AddPath(path)
}

// Open opens the folder that holds the new script, so the user can edit it.
//
// The diagnostic and the command's own output are still written the way the
// open operation writes them.
func (s createSession) Open(path string) error {
	command, output, err := s.opener.Open(path)

	if command != "" {
		s.cfg.Logger.Info("opening: ", command)
	}

	if output != "" {
		s.cfg.Logger.Out(output)
	}

	return err
}

// CreateCommandNote is the history note for a created command.
func (s createSession) CreateCommandNote(command catalog.Command) string {
	return notes.CreateCommandNote(command)
}

// CreateScriptNote is the history note for `create --script`.
func (s createSession) CreateScriptNote(command catalog.Command) string {
	return notes.CreateScriptNote(command)
}

// CreateLastNote is the history note for `create --execute`.
func (s createSession) CreateLastNote(command catalog.Command) string {
	return notes.CreateLastNote(command)
}

// CreateProcessNote is the history note for `create --process`.
func (s createSession) CreateProcessNote(command catalog.Command, args []string) string {
	return notes.CreateProcessNote(command, args)
}

// CreateSystemVariableNote is the history note for `create --system`.
func (s createSession) CreateSystemVariableNote(args []string) string {
	return notes.CreateSystemVariableNote(args)
}

// CreateSystemPathNote is the history note for `create --syspath`.
func (s createSession) CreateSystemPathNote(args []string) string {
	return notes.CreateSystemPathNote(args)
}

// CreateFileNote is the history note for `create --file`.
func (s createSession) CreateFileNote(args []string) string {
	return notes.CreateFileNote(args)
}

// CreateFolderNote is the history note for `create --folder`.
func (s createSession) CreateFolderNote(args []string) string {
	return notes.CreateFolderNote(args)
}

func init() {
	AddAllAliasFlag(createCmd, &allAlias)

	createCmd.Flags().StringVarP(&processName, "process", "p", "", "new command name for a process; pass the process name or pid as an argument")
	createCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	createCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	createCmd.Flags().BoolVar(&systemFlag, "system", false, "create system variable")
	createCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "create path variable")
	createCmd.Flags().BoolVar(&executeFlag, "execute", false, "create based on last executed command")
	createCmd.Flags().BoolVar(&scriptFlag, "script", false, "create a script file to help user create a complex command")
}
