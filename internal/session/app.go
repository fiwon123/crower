// Package session is the composition root of a crower run.
//
// It owns the configuration paths, the loaded command and history state, and
// the save/commit/revert calls that persist them. cmd builds one per run and
// hands it to the use cases through the ports they declare; nothing else may
// import this package.
package session

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/fiwon123/crower/internal/shell"
	"github.com/fiwon123/crower/pkg/logger"
)

// Session is everything one crower run needs: where its state lives, what it
// reads, and where its output goes.
type Session struct {
	History            timeline.Log
	Logger             logger.Logger
	orderKeys          []string
	allCommandsByAlias command.MapData
	allCommandsByName  command.MapData

	// Where the run's state lives. They are private because a domain should
	// not know the storage layout: it records through Save/SaveCommands and
	// reads the paths through ConfigFile/ConfigFolder. The formats themselves
	// belong to internal/adapters/storage.
	commands     *storage.CommandFile
	historyFiles *storage.HistoryFiles

	// NoInput reports that prompts must fail rather than read stdin.
	NoInput bool

	// Format is how an operation with a structured result renders it.
	Format Format

	in  *bufio.Reader
	out io.Writer
	err io.Writer
}

// ConfigFile is the configuration file this run uses.
//
// It is read-only on purpose: writing it is this package's job (SaveCommands,
// RevertTo), so a domain does not have to know the format or the location.
func (s *Session) ConfigFile() string {
	return s.commands.Path()
}

// ConfigFolder is the folder the configuration lives in, which is also where
// `create --script` keeps its scripts.
func (s *Session) ConfigFolder() string {
	return s.commands.Folder()
}

// New builds the run from an already-read configuration.
//
// allAliases must have been built by getAliasMap from allCommands; New warns
// about collisions but does not rebuild the index.
//
// It returns an error instead of panicking: a state directory that cannot be
// created, a history file that cannot be parsed, and a log that cannot be
// opened are conditions of the environment, not reasons to crash. The old code
// panicked on all three, which a script sees as an unexplained exit 2.
func New(
	cfgFilePath string,
	orderKeys []string,
	allAliases command.MapData,
	allCommands command.MapData,
	opts Options,
) (*Session, error) {
	opts = opts.withDefaults()

	commands := storage.NewCommandFile(cfgFilePath)
	historyFiles := storage.NewHistoryFiles(commands.Folder())

	if err := historyFiles.Ensure(); err != nil {
		return nil, err
	}

	logPath := filepath.Join(commands.Folder(), "crower.log")

	historyData, err := historyFiles.Load()
	if err != nil {
		return nil, err
	}

	// Entries written before the name was stored on the entry itself carry it
	// only inside the note. Resolve it on load so `--last` and the history view
	// work on a file from an older build.
	for i := range historyData.AllData {
		if historyData.AllData[i].CommandName == "" {
			historyData.AllData[i].CommandName = notes.CommandName(historyData.AllData[i].Note)
		}
	}

	log, err := logger.New(logPath, opts.Streams.Out, opts.Streams.Err)
	if err != nil {
		return nil, err
	}

	cfg := &Session{
		History:            historyData,
		commands:           commands,
		historyFiles:       historyFiles,
		Logger:             *log,
		orderKeys:          orderKeys,
		allCommandsByAlias: allAliases,
		allCommandsByName:  allCommands,
		NoInput:            opts.NoInput,
		Format:             opts.Format,
		in:                 bufio.NewReader(opts.Streams.In),
		out:                opts.Streams.Out,
		err:                opts.Streams.Err,
	}

	// A config written by an older build (or edited by hand) can hold a
	// collision. Resolution is deterministic now, but the user still needs to
	// know that one of the two commands is unreachable by that key.
	for _, warning := range commandCollisions(allCommands) {
		cfg.Logger.Warning(warning)
	}

	// A file written by an older build (or edited by hand) may hold a literal
	// the write path now refuses. Load it and say so: refusing the whole
	// configuration would leave the user unable to fix it.
	for _, command := range cfg.Catalog().Commands() {
		if err := shell.ValidateEnv(command.Env); err != nil {
			cfg.Logger.Warning(fmt.Sprintf("command %q: %v", command.Name, err))
		}
	}

	return cfg, nil
}
