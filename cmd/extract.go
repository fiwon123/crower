package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/extractcontent"

	"github.com/spf13/cobra"
)

var outDirFlag string

// Cmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract <archive...>",
	Short: "Extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...).",
	Long: `Extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...).

By default the content goes to the archive's own folder; --out changes it.

	crower extract files.zip
	crower extract files.zip --out ./build
	crower extract "*.zip"`,
	Example: `	crower extract files.zip
	crower extract files.zip --out ./build
	crower extract "*.zip" --out ./build`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "compressed file paths")
		}

		session := newExtractSession(cfg)

		return extractcontent.New(session, session, session, session).Run(
			extractcontent.Request{Args: args, OutDir: outDirFlag},
		)
	},
}

// extractSession adapts the configuration to the ports the extract use case
// declares.
type extractSession struct {
	legacySession

	archives *system.Archiver
}

// newExtractSession wires the configuration and the platform archiver.
func newExtractSession(cfg *app.Config) extractSession {
	return extractSession{
		legacySession: newLegacySession(cfg),
		archives:      system.NewArchiver(),
	}
}

// Extract expands one archive.
func (s extractSession) Extract(ext, archive, outDir string) (string, error) {
	return s.archives.Extract(ext, archive, outDir)
}

// ExtractNote is the history note for the extraction.
func (s extractSession) ExtractNote(out string, args []string) string {
	return notes.ExtractNote(out, args)
}

func init() {
	extractCmd.Flags().StringVarP(&outDirFlag, "out", "o", "", "out folder where to be extracted")
}
