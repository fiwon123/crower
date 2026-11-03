package searchdom

import (
	"testing"

	"github.com/fiwon123/crower/internal/data/operations"
)

func TestSearchPayload(t *testing.T) {
	tests := []struct {
		name string
		out  string
		want []string
	}{
		{"no match", "", []string{}},
		{"one path", "./a.txt\n", []string{"./a.txt"}},
		{"several paths", "./a.txt\n./b.txt\n", []string{"./a.txt", "./b.txt"}},
		{"carriage returns", "./a.txt\r\n./b.txt\r\n", []string{"./a.txt", "./b.txt"}},
		{"blank lines are formatting", "./a.txt\n\n./b.txt\n", []string{"./a.txt", "./b.txt"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := searchPayload(operations.File, "a.txt", ".", tt.out).Search.Results

			if got == nil {
				t.Fatal("results encoded a null list")
			}

			if len(got) != len(tt.want) {
				t.Fatalf("results = %v, want %v", got, tt.want)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("results[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSearchPayloadCarriesTheQuery(t *testing.T) {
	payload := searchPayload(operations.FileAndFolder, "report", "./docs", "")

	if payload.Search.Kind != string(operations.FileAndFolder) {
		t.Errorf("kind = %q", payload.Search.Kind)
	}

	if payload.Search.Query != "report" || payload.Search.Path != "./docs" {
		t.Errorf("query/path = %q/%q", payload.Search.Query, payload.Search.Path)
	}
}

// A browser search returns nothing on stdout: the object still has to describe
// what was asked for, and still has to carry an empty collection.
func TestSearchPayloadBrowser(t *testing.T) {
	payload := searchPayload(operations.Browser, "golang", "", "")

	if payload.Search.Kind != string(operations.Browser) {
		t.Errorf("kind = %q", payload.Search.Kind)
	}

	if payload.Search.Results == nil {
		t.Error("browser search encoded a null results list")
	}

	if len(payload.Search.Results) != 0 {
		t.Errorf("browser search returned results: %v", payload.Search.Results)
	}
}
