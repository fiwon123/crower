package main

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The Makefile declares its targets in a .PHONY list. When a target is listed
// there but its recipe is deleted, `make` treats it as an up-to-date target
// with no commands: it prints "Nothing to be done" and exits 0. That happened
// to `check` (and the whole dev-sandbox layer) in issue #90, and it meant a
// gate that reported success while running nothing.
//
// These tests are the guard: every declared target must have a rule, and the
// gate must still run gofmt, vet, test and build.

// phonyTargets returns the target names declared in the file's .PHONY list,
// following backslash continuations.
func phonyTargets(text string) []string {
	lines := strings.Split(text, "\n")
	var names []string

	for i, line := range lines {
		if !strings.HasPrefix(line, ".PHONY:") {
			continue
		}

		// The declaration may span several lines; keep consuming while the
		// current line ends with a continuation.
		joined := strings.TrimSuffix(line, "\\")
		for j := i; strings.HasSuffix(strings.TrimRight(lines[j], " \t"), "\\"); j++ {
			joined += " " + strings.TrimSuffix(strings.TrimRight(lines[j+1], " \t"), "\\")
		}

		joined = strings.TrimPrefix(joined, ".PHONY:")
		names = append(names, strings.Fields(joined)...)
	}

	return names
}

// missingRules returns the declared targets that have no rule in the file.
func missingRules(text string) []string {
	var missing []string

	for _, name := range phonyTargets(text) {
		// A rule starts at column 0 with "name:" — a recipe line is indented,
		// so it cannot match.
		if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(name) + `:`).MatchString(text) {
			missing = append(missing, name)
		}
	}

	return missing
}

func TestMakefileDeclaredTargetsHaveRules(t *testing.T) {
	data, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}

	if missing := missingRules(string(data)); len(missing) > 0 {
		t.Fatalf("declared in .PHONY but has no recipe (make would silently do nothing): %s",
			strings.Join(missing, ", "))
	}
}

func TestMakefileCheckTargetRunsTheGate(t *testing.T) {
	data, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("reading Makefile: %v", err)
	}

	text := string(data)
	if missing := missingRules(text); contains(missing, "check") {
		t.Fatalf("the check target has no recipe — `make check` would report success while doing nothing")
	}

	// All four gate steps must still be present, in order.
	checkBody := ruleBody(text, "check")
	for _, want := range []string{"gofmt -l", "go vet ./...", "go test ./...", "go build ./..."} {
		if !strings.Contains(checkBody, want) {
			t.Errorf("check recipe is missing %q; the gate no longer covers it", want)
		}
	}
}

// ruleBody returns the recipe text of the named target, up to the next
// top-level construct.
func ruleBody(text, target string) string {
	lines := strings.Split(text, "\n")
	start := -1

	for i, line := range lines {
		if strings.HasPrefix(line, target+":") {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}

	var body []string
	for _, line := range lines[start:] {
		if line != "" && !strings.HasPrefix(line, "\t") {
			break
		}
		body = append(body, line)
	}

	return strings.Join(body, "\n")
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}

	return false
}

// TestMissingRulesCatchesTheRegression proves the guard itself can fail: it is
// given the exact shape that broke in issue #90 — a target declared in .PHONY
// with its recipe removed — and must report it.
func TestMissingRulesCatchesTheRegression(t *testing.T) {
	broken := strings.Join([]string{
		".PHONY: clean release check \\",
		"\tdev-up dev-shell",
		"",
		"clean:",
		"\trm -rf build",
		"",
		"release: clean",
		"\t@echo done",
		"",
	}, "\n")

	if got := missingRules(broken); !contains(got, "check") {
		t.Errorf("missingRules() = %v, want it to include \"check\"", got)
	}
	if got := missingRules(broken); !contains(got, "dev-up") {
		t.Errorf("missingRules() = %v, want it to include \"dev-up\"", got)
	}
	if got := missingRules(broken); contains(got, "clean") || contains(got, "release") {
		t.Errorf("missingRules() = %v, want defined targets excluded", got)
	}
}
