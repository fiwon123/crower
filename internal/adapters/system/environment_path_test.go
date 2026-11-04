package system

import "testing"

// A line that is not a PATH export is not one: "export PATHDIR=…" starts with
// "export PATH" and holds a different variable.
func TestIsPathExport(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"export PATH=\"$PATH:/opt\"", true},
		{"export PATH=/usr/local/bin:$PATH", true},
		{"  export PATH=$PATH", true},
		{"export PATH", true},
		{"export PATHDIR=/opt", false},
		{"export MYPATH=/opt", false},
		{"PATH=/opt", false},
		{"# export PATH=/opt", false},
	}

	for _, tc := range cases {
		if got := isPathExport(tc.line); got != tc.want {
			t.Errorf("isPathExport(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

// The elements of the PATH are compared whole.
func TestPathHasValue(t *testing.T) {
	cases := []struct {
		line  string
		value string
		want  bool
	}{
		{"export PATH=\"$PATH:/opt/tools\"", "/opt/tools", true},
		{"export PATH=\"$PATH:/opt/tools-extra\"", "/opt/tools", false},
		{"export PATH=/opt/tools:$PATH", "/opt/tools", true},
		{"export PATH=/opt/tools-extra:$PATH", "/opt/tools", false},
		{"export PATH=\"$PATH:/opt/tools\"", "/opt", false},
	}

	for _, tc := range cases {
		if got := pathHasValue(tc.line, tc.value); got != tc.want {
			t.Errorf("pathHasValue(%q, %q) = %v, want %v", tc.line, tc.value, got, tc.want)
		}
	}
}

func TestAddValueToPathExportKeepsTheLineForm(t *testing.T) {
	cases := []struct {
		line  string
		value string
		want  string
	}{
		{"export PATH=\"$PATH\"", "/opt", "export PATH=\"$PATH:/opt\""},
		{"export PATH=$PATH", "/opt", "export PATH=$PATH:/opt"},
		{"export PATH", "/opt", "export PATH=\"$PATH:/opt\""},
	}

	for _, tc := range cases {
		if got := addValueToPathExport(tc.line, tc.value); got != tc.want {
			t.Errorf("addValueToPathExport(%q, %q) = %q, want %q", tc.line, tc.value, got, tc.want)
		}
	}
}
