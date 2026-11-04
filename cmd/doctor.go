package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/storage"

	"github.com/spf13/cobra"
)

// doctorReport is what `crower doctor` answers, in text or JSON.
type doctorReport struct {
	Setup          string `json:"setup"`
	ConfigFile     string `json:"config_file"`
	ConfigFolder   string `json:"config_folder"`
	ConfigWritable bool   `json:"config_writable"`
	Commands       int    `json:"commands"`
	History        int    `json:"history"`
	UpdateCheck    string `json:"update_check"`
	LastCheck      string `json:"last_check,omitempty"`
	Latest         string `json:"latest,omitempty"`
	Shell          string `json:"shell"`
	Editor         string `json:"editor"`
	Terminal       bool   `json:"terminal"`
	GOOS           string `json:"goos"`
	GOARCH         string `json:"goarch"`
}

var doctorJSON bool

// doctorCmd prints what a support conversation would ask for first.
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Show paths, permissions and environment crower depends on.",
	Long: `Show the configuration paths, whether they are writable, how many commands
and history entries exist, and the environment crower uses (shell, editor,
terminal). --json prints the same as one object.`,
	Example: "  crower doctor\n  crower doctor --json",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		setup := activeSetup()
		if setup == "" {
			setup = "default"
		}

		report := buildDoctorReport(setup, cfg.ConfigFile(), cfg.ConfigFolder(), cfg.Catalog().Len(), len(cfg.History.AllData))

		if doctorJSON {
			line, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(streams.Out, string(line))
			return nil
		}

		printDoctorReport(report)
		return nil
	},
}

func buildDoctorReport(setup, cfgFile, cfgFolder string, commands, history int) doctorReport {
	settings := storage.NewSettings(cfgFolder).Load()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = os.Getenv("ComSpec")
	}

	return doctorReport{
		Setup:          setup,
		ConfigFile:     cfgFile,
		ConfigFolder:   cfgFolder,
		ConfigWritable: folderWritable(cfgFolder),
		Commands:       commands,
		History:        history,
		UpdateCheck:    updateCheckLabel(settings.UpdateCheck),
		LastCheck:      settings.LastCheck,
		Latest:         settings.Latest,
		Shell:          shell,
		Editor:         os.Getenv("EDITOR"),
		Terminal:       stderrIsTerminal(),
		GOOS:           runtime.GOOS,
		GOARCH:         runtime.GOARCH,
	}
}

// updateCheckLabel names the mode, including the empty default.
func updateCheckLabel(mode string) string {
	if mode == "" {
		return "daily (default)"
	}

	return mode
}

// folderWritable reports whether crower can create its next file in folder.
func folderWritable(folder string) bool {
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return false
	}

	probe, err := os.CreateTemp(folder, ".doctor-*")
	if err != nil {
		return false
	}

	name := probe.Name()
	probe.Close()
	os.Remove(name)

	return true
}

func printDoctorReport(r doctorReport) {
	lines := []struct{ label, value string }{
		{"Setup", r.Setup},
		{"Config file", r.ConfigFile},
		{"Config folder", r.ConfigFolder},
		{"Config writable", yesNo(r.ConfigWritable)},
		{"Commands", fmt.Sprintf("%d", r.Commands)},
		{"History entries", fmt.Sprintf("%d", r.History)},
		{"Update check", r.UpdateCheck},
	}

	if r.LastCheck != "" {
		lines = append(lines, struct{ label, value string }{"Last check", r.LastCheck})
	}

	if r.Latest != "" {
		lines = append(lines, struct{ label, value string }{"Latest seen", r.Latest})
	}

	lines = append(lines,
		struct{ label, value string }{"Shell", emptyDash(r.Shell)},
		struct{ label, value string }{"Editor", emptyDash(r.Editor)},
		struct{ label, value string }{"Terminal", yesNo(r.Terminal)},
		struct{ label, value string }{"Platform", r.GOOS + "/" + r.GOARCH},
	)

	for _, line := range lines {
		fmt.Fprintf(streams.Out, "%-16s %s\n", line.label+":", line.value)
	}
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}

	return value
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}

	return "no"
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "print the report as one JSON object")
	rootCmd.AddCommand(doctorCmd)
}
