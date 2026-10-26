package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The release workflow names a runner label for every row of its build matrix.
// A label GitHub has retired is not rejected: the job queues with zero steps and
// never starts, and because `publish` needs every build, the dispatch creates no
// release at all — it just hangs. That is how both v0.5.0 dispatches were lost:
// the darwin/amd64 row sat on `macos-13`, retired in December 2025 (#121, runs
// 40343392695 and 40368033298).
//
// These tests are the guard: every matrix row must name a runner label that
// exists, and every platform the release notes offer must still be built.

const releaseWorkflowPath = ".github/workflows/release.yml"

// releaseTarget is one row of the build matrix.
type releaseTarget struct {
	goos   string
	goarch string
	runner string
}

// supportedRunners lists the runner labels the release workflow may use. The
// `*-latest` aliases belong here because GitHub keeps them on a supported
// image; a versioned label stays only while that image is supported. Retired
// labels are deliberately absent — macos-13 (retired 2025-12), macos-14
// (unsupported from 2026-11-02) and ubuntu-20.04 — so re-introducing one fails
// the suite instead of stalling the next release.
var supportedRunners = map[string]bool{
	"ubuntu-latest":  true,
	"ubuntu-24.04":   true,
	"ubuntu-22.04":   true,
	"windows-latest": true,
	"windows-2025":   true,
	"windows-2022":   true,
	"macos-latest":   true,
	"macos-15":       true,
	"macos-15-intel": true,
}

// releaseTargets returns the build matrix rows of a release workflow: the
// goos/goarch/runner triple each `- goos:` entry declares. Fields are read only
// while they are indented deeper than the row that opened them, so a `runner:`
// key belonging to another job is never attributed to a row.
func releaseTargets(workflow string) []releaseTarget {
	var (
		targets   []releaseTarget
		current   *releaseTarget
		rowIndent int
	)

	flush := func() {
		if current != nil {
			targets = append(targets, *current)
			current = nil
		}
	}

	for _, line := range strings.Split(workflow, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))

		switch {
		case strings.HasPrefix(trimmed, "- goos:"):
			flush()
			current = &releaseTarget{goos: strings.TrimSpace(strings.TrimPrefix(trimmed, "- goos:"))}
			rowIndent = indent
		case current == nil:
			// Not inside a matrix row.
		case indent <= rowIndent:
			// A key at the row's own indentation belongs to the enclosing
			// mapping: the row is over.
			flush()
		case strings.HasPrefix(trimmed, "goarch:"):
			current.goarch = strings.TrimSpace(strings.TrimPrefix(trimmed, "goarch:"))
		case strings.HasPrefix(trimmed, "runner:"):
			current.runner = strings.TrimSpace(strings.TrimPrefix(trimmed, "runner:"))
		}
	}
	flush()

	return targets
}

func TestReleaseWorkflowBuildsEveryAdvertisedPlatform(t *testing.T) {
	targets := releaseTargets(readRepositoryFile(t, releaseWorkflowPath))
	if len(targets) == 0 {
		t.Fatal("no build matrix rows found in the release workflow — the guard is not looking at anything")
	}

	built := make(map[string]bool, len(targets))
	for _, target := range targets {
		built[target.goos+"/"+target.goarch] = true
	}

	// The release notes offer a download for each of these.
	for _, platform := range []string{
		"linux/amd64", "linux/arm64",
		"windows/amd64", "windows/arm64",
		"darwin/amd64", "darwin/arm64",
	} {
		if !built[platform] {
			t.Errorf("the release no longer builds %s, but its notes still offer that download", platform)
		}
	}
}

func TestReleaseWorkflowRunnersAreSupported(t *testing.T) {
	targets := releaseTargets(readRepositoryFile(t, releaseWorkflowPath))
	if len(targets) == 0 {
		t.Fatal("no build matrix rows found in the release workflow — the guard is not looking at anything")
	}

	for _, target := range targets {
		platform := target.goos + "/" + target.goarch
		if target.goos == "" || target.goarch == "" {
			t.Errorf("matrix row %+v has no goos/goarch, so its job name and artifacts are undefined", target)
			continue
		}
		if target.runner == "" {
			t.Errorf("%s names no runner, so the job cannot be scheduled", platform)
			continue
		}
		if !supportedRunners[target.runner] {
			t.Errorf("%s builds on %q, which is not a runner label known to exist — a retired label "+
				"queues forever instead of failing, so the dispatch never reaches the publish job (#121)",
				platform, target.runner)
		}
	}
}

// TestReleaseTargetsCatchRetiredRunner proves the guard can fail: it is given
// the exact matrix the workflow shipped with in issue #121 — darwin/amd64 on
// macos-13 — and must report it.
func TestReleaseTargetsCatchRetiredRunner(t *testing.T) {
	shipped := strings.Join([]string{
		"    strategy:",
		"      fail-fast: false",
		"      matrix:",
		"        include:",
		"          - goos: linux",
		"            goarch: amd64",
		"            runner: ubuntu-latest",
		"          - goos: darwin",
		"            goarch: amd64",
		"            runner: macos-13",
		"    steps:",
		"      - name: Checkout tag",
		"        with:",
		"          runner: not-a-matrix-row",
		"",
	}, "\n")

	targets := releaseTargets(shipped)
	if len(targets) != 2 {
		t.Fatalf("releaseTargets() returned %d rows, want 2 — it is not parsing the matrix shape", len(targets))
	}

	if targets[0].runner != "ubuntu-latest" || targets[1].runner != "macos-13" {
		t.Fatalf("releaseTargets() = %+v, want the two rows in order with their runners", targets)
	}

	if supportedRunners[targets[1].runner] {
		t.Errorf("supportedRunners[%q] = true, want false: the label was retired in December 2025", targets[1].runner)
	}

	if targets[0].goarch != "amd64" || targets[1].goarch != "amd64" {
		t.Errorf("releaseTargets() = %+v, want each row to keep its goarch", targets)
	}
}

// --- guards for every workflow file ----------------------------------------
//
// Two faults in the release notes step cost the v0.5.0 dispatch its draft
// (#124):
//
//   - the step ran `gh`, which authenticates from the environment, without
//     `GITHUB_TOKEN` in its `env:` — so the call failed and `2>/dev/null` hid
//     the reason;
//   - its fallback wrote a trailing underscore after the expansion, which bash
//     reads as part of the variable name, so `set -u` aborted the job.
//
// Both are mechanical, so both are checked mechanically.

// workflowStep is one `- name:` entry of a job, with the lines that follow it.
type workflowStep struct {
	name string
	body string
}

// workflowSteps returns the steps of every job in a workflow. A step ends at
// the next `- name:` or at a key indented to the job level or above, so a
// `runner:` or `steps:` belonging to another job is never read as part of it.
func workflowSteps(workflow string) []workflowStep {
	var (
		steps   []workflowStep
		current *workflowStep
	)

	flush := func() {
		if current != nil {
			steps = append(steps, *current)
			current = nil
		}
	}

	for _, line := range strings.Split(workflow, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))

		switch {
		case strings.HasPrefix(trimmed, "- name:"):
			flush()
			current = &workflowStep{name: strings.TrimSpace(strings.TrimPrefix(trimmed, "- name:"))}
		case current == nil:
			// Not inside a step.
		case indent <= 4:
			// `jobs:`, `  publish:` and `    runs-on:` are job-level keys: the
			// step above them is over.
			flush()
		default:
			current.body += line + "\n"
		}
	}
	flush()

	return steps
}

var ghCommand = regexp.MustCompile(`\bgh [a-z]+`)

// stepsCallingGh returns the steps that invoke the GitHub CLI. A token set at
// workflow level would satisfy `gh` too, but no workflow here does that, so a
// per-step `GITHUB_TOKEN` is the rule this guard enforces.
func stepsCallingGh(workflow string) []workflowStep {
	var calling []workflowStep

	for _, step := range workflowSteps(workflow) {
		if ghCommand.MatchString(step.body) {
			calling = append(calling, step)
		}
	}

	return calling
}

// trailingUnderscoreExpansion matches an expansion whose variable name ends in
// an underscore because the author meant braces around the name. `$NAME_X` is a
// different variable and does not match, which is what keeps `$GITHUB_OUTPUT`
// out of the results.
var trailingUnderscoreExpansion = regexp.MustCompile(`\$[A-Za-z_][A-Za-z0-9_]*_([^A-Za-z0-9_]|$)`)

// ambiguousExpansions returns the trailing-underscore expansions in a script.
// Shell comments are skipped: they do not execute, so they cannot be unbound.
func ambiguousExpansions(script string) []string {
	var found []string

	for i, line := range strings.Split(script, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if match := trailingUnderscoreExpansion.FindString(line); match != "" {
			found = append(found, fmt.Sprintf("line %d: %s", i+1, strings.TrimSpace(match)))
		}
	}

	return found
}

func workflowFilesToGuard(t *testing.T) []string {
	t.Helper()

	files, err := filepath.Glob(".github/workflows/*.yml")
	if err != nil {
		t.Fatalf("globbing the workflows: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no workflow files found — the guard is not looking at anything")
	}

	return files
}

func TestWorkflowStepsCallingGhSetToken(t *testing.T) {
	sawGh := false

	for _, file := range workflowFilesToGuard(t) {
		for _, step := range stepsCallingGh(readRepositoryFile(t, file)) {
			sawGh = true
			if !strings.Contains(step.body, "GITHUB_TOKEN") {
				t.Errorf("%s → step %q runs gh without GITHUB_TOKEN in its env: gh authenticates from the "+
					"environment, so the call fails there and whatever fallback the step carries stands in for "+
					"the real result (#124)", file, step.name)
			}
		}
	}

	if !sawGh {
		t.Fatal("no step invokes gh — the guard is not looking at anything")
	}
}

func TestWorkflowExpansionsAreNotAmbiguous(t *testing.T) {
	for _, file := range workflowFilesToGuard(t) {
		if found := ambiguousExpansions(readRepositoryFile(t, file)); len(found) > 0 {
			t.Errorf("%s ends an expansion with an underscore where braces were meant: bash reads that "+
				"underscore as part of the variable name, and the unset variable aborts the step under "+
				"`set -u` (#124):\n  %s", file, strings.Join(found, "\n  "))
		}
	}
}

// TestWorkflowGuardsCatchTheReleaseNotesBug proves both guards can fail: given
// the step as it shipped in the v0.5.0 dispatch, they must report it.
func TestWorkflowGuardsCatchTheReleaseNotesBug(t *testing.T) {
	shipped := strings.Join([]string{
		"jobs:",
		"  publish:",
		"    runs-on: ubuntu-latest",
		"    steps:",
		"      - name: Create professional release notes",
		"        id: notes",
		"        run: |",
		"          set -euo pipefail",
		`          PREV="v0.4.0"`,
		`          gh api repos/x/y/releases/generate-notes --jq .body >> notes.md 2>/dev/null || echo "_No changes since $PREV_" >> notes.md`,
		"      - name: Create or refresh draft release",
		"        env:",
		"          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}",
		"        run: |",
		`          gh release create v0.5.0 --draft dist/*`,
		"",
	}, "\n")

	steps := stepsCallingGh(shipped)
	if len(steps) != 2 {
		t.Fatalf("stepsCallingGh() found %d steps, want 2 — the guard is not parsing the step shape", len(steps))
	}

	if strings.Contains(steps[0].body, "GITHUB_TOKEN") {
		t.Errorf("the notes step reports a token, but it shipped without one")
	}

	found := ambiguousExpansions(steps[0].body)
	if len(found) != 1 || !strings.Contains(found[0], "$PREV_") {
		t.Errorf("ambiguousExpansions() = %v, want exactly the shipped expansion", found)
	}

	if !strings.Contains(steps[1].body, "GITHUB_TOKEN") {
		t.Errorf("the draft step reports no token, but it ships with one — the guard would flag it wrongly")
	}

	// The corrected shape must pass the expansion guard, or it would be unusable.
	corrected := strings.Replace(shipped, `$PREV_"`, `${PREV}_"`, 1)
	if found := ambiguousExpansions(corrected); len(found) > 0 {
		t.Errorf("ambiguousExpansions() = %v, want none once braces close the name", found)
	}

	// `$GITHUB_OUTPUT`-style names contain an underscore but are unambiguous.
	if found := ambiguousExpansions(`echo "commit=$MAIN_COMMIT" >> "$GITHUB_OUTPUT"`); len(found) > 0 {
		t.Errorf("ambiguousExpansions() = %v, want none: an internal underscore is not a trailing one", found)
	}
}
