package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/update"
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
// releaseJobSection returns one top-level job from a workflow, from its
// two-space key to the next job key. The release has two matrices now (the
// cgo-free build and the Fyne build); a parser that scans the whole file would
// mix them.
func releaseJobSection(workflow, job string) string {
	lines := strings.Split(workflow, "\n")
	start := -1

	for i, line := range lines {
		if line == "  "+job+":" {
			start = i
			continue
		}

		if start < 0 || i == start {
			continue
		}

		if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") &&
			!strings.HasPrefix(strings.TrimSpace(line), "#") &&
			strings.HasSuffix(strings.TrimSpace(line), ":") {
			return strings.Join(lines[start:i], "\n")
		}
	}

	if start >= 0 {
		return strings.Join(lines[start:], "\n")
	}

	return ""
}

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

// releaseTargetsForJobs merges the matrix rows of several release jobs. The
// release builds four platforms natively (gui) and two arm64 platforms with
// fyne-cross (gui-cross), so no single job lists all six.
func releaseTargetsForJobs(workflow string, jobs ...string) []releaseTarget {
	var all []releaseTarget
	for _, job := range jobs {
		all = append(all, releaseTargets(releaseJobSection(workflow, job))...)
	}
	return all
}

func TestReleaseWorkflowBuildsEveryAdvertisedPlatform(t *testing.T) {
	targets := releaseTargetsForJobs(readRepositoryFile(t, releaseWorkflowPath), "gui", "gui-cross")
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

// The native GUI job is a CGO build. It targets Linux/Windows amd64 and macOS
// arm64/Intel; the two arm64 targets are built by gui-cross.
func TestReleaseWorkflowBuildsTheGUIPlatforms(t *testing.T) {
	targets := releaseTargets(releaseJobSection(readRepositoryFile(t, releaseWorkflowPath), "gui"))
	if len(targets) != 4 {
		t.Fatalf("found %d GUI targets, want 4", len(targets))
	}

	built := make(map[string]bool, len(targets))
	for _, target := range targets {
		built[target.goos+"/"+target.goarch] = true

		if !supportedRunners[target.runner] {
			t.Errorf("GUI %s/%s builds on %q, which is not a runner label known to exist", target.goos, target.goarch, target.runner)
		}
	}

	for _, platform := range []string{"linux/amd64", "windows/amd64", "darwin/arm64", "darwin/amd64"} {
		if !built[platform] {
			t.Errorf("the GUI release no longer builds %s", platform)
		}
	}
}

// The arm64 GUI binaries are built with fyne-cross in Docker on ubuntu-latest.
// They must stay on a runner that has Docker, and the archive name must match
// the platform the CLI pair is built for.
func TestReleaseWorkflowBuildsTheGUICrossPlatforms(t *testing.T) {
	targets := releaseTargets(releaseJobSection(readRepositoryFile(t, releaseWorkflowPath), "gui-cross"))
	if len(targets) != 2 {
		t.Fatalf("found %d cross GUI targets, want 2", len(targets))
	}

	built := make(map[string]bool, len(targets))
	for _, target := range targets {
		built[target.goos+"/"+target.goarch] = true

		if target.runner != "ubuntu-latest" {
			t.Errorf("cross GUI %s/%s builds on %q, want ubuntu-latest (fyne-cross needs Docker)", target.goos, target.goarch, target.runner)
		}
	}

	for _, platform := range []string{"linux/arm64", "windows/arm64"} {
		if !built[platform] {
			t.Errorf("the cross GUI release no longer builds %s", platform)
		}
	}
}

func TestReleaseWorkflowRunnersAreSupported(t *testing.T) {
	targets := releaseTargetsForJobs(readRepositoryFile(t, releaseWorkflowPath), "gui", "gui-cross")
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

// The upgrade downloads one asset per release target. Its name has to carry the
// platform and the architecture the build matrix actually produces, or the
// request is a 404 (#143).
func TestEveryReleaseTargetHasAnUpgradeAsset(t *testing.T) {
	targets := releaseTargetsForJobs(readRepositoryFile(t, releaseWorkflowPath), "gui", "gui-cross")
	if len(targets) != 6 {
		t.Fatalf("found %d release targets, want 6", len(targets))
	}

	// The Downloads table names the archives the build emits. Substituting a
	// sample tag lets the asset name be compared with both sides: the matrix
	// and the filenames the release actually publishes.
	downloads := strings.ReplaceAll(readRepositoryFile(t, releaseWorkflowPath), "${TAG}", "v1.2.3")
	downloads = strings.ReplaceAll(downloads, "@TAG@", "v1.2.3")

	for _, target := range targets {
		asset, err := update.AssetName(target.goos, target.goarch, "v1.2.3")
		if err != nil {
			t.Errorf("%s/%s: AssetName error = %v", target.goos, target.goarch, err)
			continue
		}

		if !strings.Contains(downloads, asset) {
			t.Errorf("%s/%s: the release does not publish %q", target.goos, target.goarch, asset)
		}
	}
}

// The published body follows the v0.8.0 house style: an opening sentence then
// section emoji in a fixed order, with platform icons in a CLI/GUI table (#353).
// Shields.io badges are not part of it — GitHub has no custom text colour, and
// the v0.8.0 page that set the style has none.
func TestReleaseNotesFollowTheHouseStyle(t *testing.T) {
	notes := releaseJobSection(readRepositoryFile(t, releaseWorkflowPath), "publish")

	for _, banned := range []string{"img.shields.io", "## Crower"} {
		if strings.Contains(notes, banned) {
			t.Errorf("the notes step still emits %q; the house style uses section emoji and platform icons, not badges or a repeated title (#353)", banned)
		}
	}

	headings := []string{
		"## 📦 Downloads",
		"## ⚙️ Install",
		"## ⚡ Upgrade",
		"## 🔐 Verify",
		"## 📝 Changelog",
	}

	last := -1
	for _, heading := range headings {
		at := strings.Index(notes, heading)
		if at < 0 {
			t.Fatalf("the notes step does not emit %q", heading)
		}
		if at <= last {
			t.Errorf("the notes step emits %q out of order", heading)
		}
		last = at
	}

	rows := []string{
		"| Platform | Archive |",
		"| 🐧 Linux amd64 |",
		"| 🐧 Linux arm64 |",
		"| 🪟 Windows amd64 |",
		"| 🪟 Windows arm64 |",
		"| 🍎 macOS Intel |",
		"| 🍎 macOS Apple Silicon |",
	}
	for _, row := range rows {
		if !strings.Contains(notes, row) {
			t.Errorf("the Downloads table is missing %q", row)
		}
	}

	if !strings.Contains(notes, "vcs.modified=false") {
		t.Errorf("the Verify section does not state vcs.modified=false")
	}
}

// The template is the contract the workflow and the hand-written bodies follow,
// so it has to keep the full order — including the two sections the workflow
// leaves for the release manager to write.
func TestReleaseNotesTemplateKeepsTheHouseOrder(t *testing.T) {
	template := readRepositoryFile(t, "docs/release-notes-template.md")

	headings := []string{
		"## ✨ Highlights",
		"## 🛠 Fixes",
		"## 📦 Downloads",
		"## ⚙️ Install",
		"## ⚡ Upgrade",
		"## 🔐 Verify",
		"## 📝 Changelog",
	}

	last := -1
	for _, heading := range headings {
		at := strings.Index(template, heading)
		if at < 0 {
			t.Fatalf("the release notes template does not show %q", heading)
		}
		if at <= last {
			t.Errorf("the release notes template shows %q out of order", heading)
		}
		last = at
	}

	if strings.Contains(template, "img.shields.io") {
		t.Errorf("the release notes template mentions shields.io; colour is emoji, not badges")
	}

	if !strings.Contains(template, "| Platform | Archive |") {
		t.Errorf("the release notes template does not show the one-archive table")
	}
}

// The style is only a guarantee if the workflow builds it and then checks it:
// `releasehighlights` writes the Highlights and Fixes sections from the merged
// PRs, and `bodycheck` refuses a body that drifts from the template.
func TestReleaseWorkflowGeneratesAndChecksTheBodyShape(t *testing.T) {
	workflow := readRepositoryFile(t, releaseWorkflowPath)

	tools := []string{
		"go run ./internal/tools/releasehighlights",
		"go run ./internal/tools/bodycheck notes.md",
	}
	for _, tool := range tools {
		if !strings.Contains(workflow, tool) {
			t.Errorf("the release workflow does not run %q; the house style is no longer enforced", tool)
		}
	}

	generated := strings.Index(workflow, tools[0])
	checked := strings.Index(workflow, tools[1])
	if generated < 0 || checked < 0 || checked < generated {
		t.Errorf("the body check must run after the sections are generated (generated at %d, checked at %d)", generated, checked)
	}

	// The notes step now calls `gh api` and `gh pr list`, so it needs the token.
	// TestWorkflowStepsCallingGhSetToken checks this globally; this assertion
	// keeps the reason next to the style guard.
	notes := releaseJobSection(workflow, "publish")
	if !strings.Contains(notes, "GITHUB_TOKEN") {
		t.Errorf("the publish job has no GITHUB_TOKEN; the PR list call would fail (#124)")
	}
}
