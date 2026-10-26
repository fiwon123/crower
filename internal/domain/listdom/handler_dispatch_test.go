package listdom

import (
	"errors"
	"runtime"
	"testing"
)

// The platform dispatch is tested through the pure command builders rather than
// through the Handler methods: the builders decide what runs, they need no
// application, and they never spawn a shell.
func TestListFolderCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		path    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "/tmp/example", "dir '/tmp/example'", false},
		{"linux", "linux", "/tmp/example", "\"ls '/tmp/example'\"", false},
		{"darwin", "darwin", "/tmp/example", "\"ls '/tmp/example'\"", false},
		{"unsupported freebsd", "freebsd", "/tmp/example", "", true},
		{"unsupported empty", "", "/tmp/example", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listFolderCommand(tt.goos, tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("listFolderCommand(%q) error = %v, wantErr %v", tt.goos, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrUnsupportedOS) {
				t.Errorf("listFolderCommand(%q) error = %v, want ErrUnsupportedOS", tt.goos, err)
			}
			if got != tt.want {
				t.Errorf("listFolderCommand(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}

func TestListSystemCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "'set'", false},
		{"linux", "linux", "'printenv'", false},
		{"darwin", "darwin", "'printenv'", false},
		{"unsupported freebsd", "freebsd", "", true},
		{"unsupported empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listSystemCommand(tt.goos)
			if (err != nil) != tt.wantErr {
				t.Fatalf("listSystemCommand(%q) error = %v, wantErr %v", tt.goos, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrUnsupportedOS) {
				t.Errorf("listSystemCommand(%q) error = %v, want ErrUnsupportedOS", tt.goos, err)
			}
			if got != tt.want {
				t.Errorf("listSystemCommand(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}

func TestListSysPathCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "'echo %PATH%'", false},
		{"linux", "linux", "'echo $PATH'", false},
		{"darwin", "darwin", "'echo $PATH'", false},
		{"unsupported freebsd", "freebsd", "", true},
		{"unsupported empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listSysPathCommand(tt.goos)
			if (err != nil) != tt.wantErr {
				t.Fatalf("listSysPathCommand(%q) error = %v, wantErr %v", tt.goos, err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrUnsupportedOS) {
				t.Errorf("listSysPathCommand(%q) error = %v, want ErrUnsupportedOS", tt.goos, err)
			}
			if got != tt.want {
				t.Errorf("listSysPathCommand(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}

// An unsupported platform must fail every listing method with an error and no
// output. execute is deliberately nil: the failure has to happen before
// anything is executed, so a handler that reached the shell would panic here
// instead of returning.
func TestUnsupportedOSReturnsError(t *testing.T) {
	handler := &Handler{goos: "plan9"}

	tests := []struct {
		name string
		list func() (string, error)
	}{
		{"ListFolder", func() (string, error) { return handler.ListFolder("/tmp/example") }},
		{"ListSystem", handler.ListSystem},
		{"ListSysPath", handler.ListSysPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.list()
			if !errors.Is(err, ErrUnsupportedOS) {
				t.Fatalf("%s error = %v, want ErrUnsupportedOS", tt.name, err)
			}
			if out != "" {
				t.Errorf("%s out = %q, want empty", tt.name, out)
			}
		})
	}
}

// NewHandler has to carry the running platform, or every command falls through
// to the unsupported arm on the user's machine.
func TestNewHandlerUsesRunningOS(t *testing.T) {
	handler := NewHandler(nil)
	if handler.goos != runtime.GOOS {
		t.Errorf("NewHandler goos = %q, want %q", handler.goos, runtime.GOOS)
	}
}
