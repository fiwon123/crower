package main

import (
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
