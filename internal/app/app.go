package app

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"

	"github.com/fiwon123/crower/pkg/logger"
	"github.com/fiwon123/crower/pkg/utils"
)

// Config is everything one crower run needs: where its state lives, what it
// reads, and where its output goes.
type Config struct {
	CfgFilePath        string
	HistoryFilePath    string
	HistoryFolderPath  string
	History            history.Data
	Logger             logger.Logger
	OrderKeys          []string
	AllCommandsByAlias command.MapData
	AllCommandsByName  command.MapData

	// NoInput reports that prompts must fail rather than read stdin.
	NoInput bool

	in  *bufio.Reader
	out io.Writer
	err io.Writer
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

	log, err := logger.New(logPath, opts.Streams.Out, opts.Streams.Err)
	if err != nil {
		return nil, err
	}

	return &Config{
		CfgFilePath:        cfgFilePath,
		History:            history,
		HistoryFilePath:    historyFilePath,
		HistoryFolderPath:  historyFolderPath,
		Logger:             *log,
		OrderKeys:          orderKeys,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
		NoInput:            opts.NoInput,
		in:                 bufio.NewReader(opts.Streams.In),
		out:                opts.Streams.Out,
		err:                opts.Streams.Err,
	}, nil
}
