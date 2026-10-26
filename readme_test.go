package main

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// The README is the only install documentation, and it is not compiled: a
// command naming a make target or an archive that does not exist survives every
// other check. These tests are the guard (issue #98):
//
//   - every `make <target>` in a README code block must be a rule in the
//     Makefile, and
//   - every archive a README code block extracts must name an architecture,
//     because a release ships both an arch-suffixed archive and a legacy
//     arch-less one, and a glob without an arch matches more than one file.

var (
	makeInvocation   = regexp.MustCompile(`\bmake\s+([a-zA-Z0-9_][a-zA-Z0-9_-]*)`)
	archiveExtract   = regexp.MustCompile(`\b(?:tar\s+-[a-zA-Z]*x[a-zA-Z]*f|unzip)\s+(.+)$`)
	architectureName = regexp.MustCompile(`(amd64|arm64)`)
)

// codeBlocks returns the contents of the fenced code blocks in a Markdown
// document. Commands live in fences; prose does not, so "make sure" in a
// sentence is not mistaken for a make invocation.
func codeBlocks(markdown string) []string {
	var blocks []string
	var current []string
	inside := false

	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			if inside {
				blocks = append(blocks, strings.Join(current, "\n"))
				current = nil
			}
			inside = !inside
			continue
		}
		if inside {
			current = append(current, line)
		}
	}

	return blocks
}

// makeTargetsInvoked returns every make target the code blocks invoke, in order
// of first appearance.
func makeTargetsInvoked(code string) []string {
	var targets []string
	seen := map[string]bool{}

	for _, match := range makeInvocation.FindAllStringSubmatch(code, -1) {
		target := match[1]
		if seen[target] {
			continue
		}
		seen[target] = true
		targets = append(targets, target)
	}

	return targets
}

// undefinedMakeTargets returns the given targets that the Makefile has no rule
// for. A rule starts at column 0 with "name:"; a recipe line is indented, so it
// cannot match.
func undefinedMakeTargets(targets []string, makefile string) []string {
	var undefined []string

	for _, target := range targets {
		if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(target) + `:`).MatchString(makefile) {
			undefined = append(undefined, target)
		}
	}

	return undefined
}

// archiveCommands returns the file argument of every archive-extracting command
// in the given code.
func archiveCommands(code string) []string {
	var commands []string

	for _, line := range strings.Split(code, "\n") {
		if match := archiveExtract.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			commands = append(commands, strings.TrimSpace(match[1]))
		}
	}

	return commands
}

// archivesWithoutArchitecture returns the archive arguments that do not name an
// architecture, and so cannot be guaranteed to match exactly one archive.
func archivesWithoutArchitecture(commands []string) []string {
	var ambiguous []string

	for _, command := range commands {
		if !architectureName.MatchString(command) {
			ambiguous = append(ambiguous, command)
		}
	}

	return ambiguous
}

func readRepositoryFile(t *testing.T, name string) string {
	t.Helper()

	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}

	return string(data)
}

func TestReadmeMakeTargetsExist(t *testing.T) {
	readme := readRepositoryFile(t, "README.md")
	makefile := readRepositoryFile(t, "Makefile")

	targets := makeTargetsInvoked(strings.Join(codeBlocks(readme), "\n"))
	if len(targets) == 0 {
		t.Fatal("no make invocations found in the README code blocks — the guard is not looking at anything")
	}

	if undefined := undefinedMakeTargets(targets, makefile); len(undefined) > 0 {
		t.Fatalf("README invokes make targets the Makefile does not define (the command fails with "+
			"\"No rule to make target\"): %s", strings.Join(undefined, ", "))
	}
}

func TestReadmeArchiveCommandsNameOneArchive(t *testing.T) {
	readme := readRepositoryFile(t, "README.md")

	commands := archiveCommands(strings.Join(codeBlocks(readme), "\n"))
	if len(commands) == 0 {
		t.Fatal("no archive extraction commands found in the README code blocks — the guard is not looking at anything")
	}

	if ambiguous := archivesWithoutArchitecture(commands); len(ambiguous) > 0 {
		t.Fatalf("README extracts archives without naming an architecture, so the glob matches more than "+
			"one file: %s", strings.Join(ambiguous, ", "))
	}
}

// TestReadmeGuardsCatchTheRegressions proves the guards can fail: they are given
// the exact shape the README shipped with in issue #98 and must report it.
func TestReadmeGuardsCatchTheRegressions(t *testing.T) {
	makefile := readRepositoryFile(t, "Makefile")

	shipped := "```bash\n" +
		"tar -xzf crower_*.tar.gz\n" +
		"chmod +x crower cr\n" +
		"```\n" +
		"```bash\n" +
		"git clone https://github.com/fiwon123/crower.git\n" +
		"cd crower\n" +
		"make build\n" +
		"sudo make install  # or add build/ to PATH\n" +
		"```\n"

	code := strings.Join(codeBlocks(shipped), "\n")

	undefined := undefinedMakeTargets(makeTargetsInvoked(code), makefile)
	for _, want := range []string{"build", "install"} {
		if !contains(undefined, want) {
			t.Errorf("undefinedMakeTargets() = %v, want it to include %q", undefined, want)
		}
	}

	if ambiguous := archivesWithoutArchitecture(archiveCommands(code)); len(ambiguous) != 1 {
		t.Errorf("archivesWithoutArchitecture() = %v, want exactly the ambiguous `crower_*.tar.gz`", ambiguous)
	}

	// A correct block must pass both guards, or the guard would be unusable.
	correct := "```bash\n" +
		"VERSION=v0.4.0\n" +
		"tar -xzf \"crower_${VERSION}_linux_amd64.tar.gz\"\n" +
		"make check\n" +
		"```\n"

	code = strings.Join(codeBlocks(correct), "\n")
	if undefined := undefinedMakeTargets(makeTargetsInvoked(code), makefile); len(undefined) > 0 {
		t.Errorf("undefinedMakeTargets() = %v, want none for a correct block", undefined)
	}
	if ambiguous := archivesWithoutArchitecture(archiveCommands(code)); len(ambiguous) > 0 {
		t.Errorf("archivesWithoutArchitecture() = %v, want none for a correct block", ambiguous)
	}
}
