package app

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/pkg/logger"
	"github.com/fiwon123/crower/pkg/utils"
)

// Config is everything one crower run needs: where its state lives, what it
// reads, and where its output goes.
type Config struct {
	History            history.Data
	Logger             logger.Logger
	OrderKeys          []string
	AllCommandsByAlias command.MapData
	AllCommandsByName  command.MapData

	// Where the run's state lives. They are private because a domain should
	// not know the storage layout: it records through Save/SaveCommands and
	// reads the paths through ConfigFile/ConfigFolder.
	cfgFilePath       string
	historyFilePath   string
	historyFolderPath string

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
func (app *Config) ConfigFile() string {
	return app.cfgFilePath
}

// ConfigFolder is the folder the configuration lives in, which is also where
// `create --script` keeps its scripts.
func (app *Config) ConfigFolder() string {
	return filepath.Dir(app.cfgFilePath)
}

// NewApp builds the run from an already-read configuration.
//
// It returns an error instead of panicking: a state directory that cannot be
// created, a history file that cannot be parsed, and a log that cannot be
// opened are conditions of the environment, not reasons to crash. The old code
// panicked on all three, which a script sees as an unexplained exit 2.
func NewApp(
	cfgFilePath string,
	orderKeys []string,
	allAliases command.MapData,
	allCommands command.MapData,
	opts Options,
) (*Config, error) {
	opts = opts.withDefaults()

	folderPath := filepath.Dir(cfgFilePath)

	historyFilePath := filepath.Join(folderPath, "historyhelper.json")
	if err := utils.CreateFileIfNotExists(historyFilePath); err != nil {
		return nil, fmt.Errorf("prepare history file %s: %w", historyFilePath, err)
	}

	historyFolderPath := filepath.Join(folderPath, "history")
	if err := utils.CreateFolderIfNotExists(historyFolderPath); err != nil {
		return nil, fmt.Errorf("prepare history folder %s: %w", historyFolderPath, err)
	}

	logPath := filepath.Join(folderPath, "crower.log")

	var history history.Data
	if err := utils.ReadJson(historyFilePath, &history); err != nil {
		return nil, fmt.Errorf("read history %s: %w", historyFilePath, err)
	}

	// Entries written before the name was stored on the entry itself carry it
	// only inside the note. Resolve it on load so `--last` and the history view
	// work on a file from an older build.
	for i := range history.AllData {
		if history.AllData[i].CommandName == "" {
			history.AllData[i].CommandName = note.CommandName(history.AllData[i].Note)
		}
	}

	log, err := logger.New(logPath, opts.Streams.Out, opts.Streams.Err)
	if err != nil {
		return nil, err
	}

	return &Config{
		History:            history,
		cfgFilePath:        cfgFilePath,
		historyFilePath:    historyFilePath,
		historyFolderPath:  historyFolderPath,
		Logger:             *log,
		OrderKeys:          orderKeys,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
		NoInput:            opts.NoInput,
		Format:             opts.Format,
		in:                 bufio.NewReader(opts.Streams.In),
		out:                opts.Streams.Out,
		err:                opts.Streams.Err,
	}, nil
}
