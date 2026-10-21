package opendom_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/opendom"
	"github.com/fiwon123/crower/internal/helper"
)

// recordingOpener returns a CommandRunner that records the command it was
// given and never executes it.
//
// The real runner shells out to xdg-open, which needs a desktop environment:
// on a headless machine it exits 3 with "no method available for opening",
// so a test that used it would report on the machine rather than on the
// code. The command is still built by the production code, so the test can
// assert on what would have run.
func recordingOpener(t *testing.T, commands *[]string) opendom.CommandRunner {
	t.Helper()

	return func(command string) (string, error) {
		*commands = append(*commands, command)
		return "", nil
	}
}

// assertOpensPath checks that the recorded command opens the given path.
func assertOpensPath(t *testing.T, commands []string, path string) {
	t.Helper()

	if len(commands) == 0 {
		t.Fatalf("no command was executed, expected one opening %q", path)
	}

	for _, command := range commands {
		if strings.Contains(command, path) {
			return
		}
	}

	t.Fatalf("no command opens %q, got %v", path, commands)
}

func TestOpenFile(t *testing.T) {
	t.Run("Open File single name", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)

		var commands []string
		openHandler := opendom.NewHandlerWithRunner(app, recordingOpener(t, &commands))

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = openHandler.Open([]string{newFilePath})
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		assertOpensPath(t, commands, newFilePath)
	})

	t.Run("Open File name with space", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)

		var commands []string
		openHandler := opendom.NewHandlerWithRunner(app, recordingOpener(t, &commands))

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = openHandler.Open([]string{newFilePath})
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		assertOpensPath(t, commands, newFilePath)
	})
}

func TestOpenFolder(t *testing.T) {
	t.Run("Open c single name", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)

		var commands []string
		openHandler := opendom.NewHandlerWithRunner(app, recordingOpener(t, &commands))

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = openHandler.Open([]string{newFolderPath})
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		assertOpensPath(t, commands, newFolderPath)
	})

	t.Run("Open Folder name with space", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)

		var commands []string
		openHandler := opendom.NewHandlerWithRunner(app, recordingOpener(t, &commands))

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = openHandler.Open([]string{newFolderPath})
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		assertOpensPath(t, commands, newFolderPath)
	})
}

func TestOpenSystem(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	openHandler := opendom.NewHandler(app)

	err = openHandler.OpenSystem()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

}
