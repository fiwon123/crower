package main

// This file is the architecture gate.
//
// The refactor in #215 moves packages between four rings -- domain,
// application, adapters and the cobra driver -- and the whole point is that the
// dependency arrow only points one way. A comment describing that rule is not
// enforcement: a single import added during a migration turns the arrow around,
// compiles, and passes every other test.
//
// The gate parses imports directly, so it needs no toolchain flag, no build
// tag, and no network. It applies every rule whose directory exists, so a rule
// for a ring that has not been created yet is inert rather than silently green,
// and a rule that is required reports when it matched nothing.

import (
	"errors"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/fiwon123/crower"

// importRule describes what one ring may not import.
type importRule struct {
	// name is what the failure message calls this ring.
	name string

	// dir is the module-relative directory the rule applies to.
	dir string

	// required fails the test when no package is found under dir, so a typo or
	// a move cannot turn the rule into a no-op that reports green.
	required bool

	// forbid lists module-relative import prefixes the ring must not reach.
	forbid []string

	// noModule forbids importing this module at all.
	noModule bool

	// noModuleOutsideOwn allows this module's imports only from the ring's own
	// subtree -- the rule that keeps a domain independent of its neighbours.
	noModuleOutsideOwn bool

	// noExternal forbids third-party imports (any path whose first element
	// contains a dot), so a ring cannot quietly take a dependency.
	noExternal bool

	// skip exempts a package from the rule. It exists so the legacy
	// `internal/domain/<op>dom` packages can stay where they are during the
	// migration without switching the rule off for the ring that replaces
	// them; #215 Phase 5 deletes them and the exemption goes with them.
	skip func(relativeDir string) bool
}

// legacyOperationPackage recognises the `<op>dom` packages the refactor
// replaces. They predate the rings, so they are exempt from the ring rules
// until they are deleted rather than moved.
func legacyOperationPackage(relativeDir string) bool {
	return strings.HasSuffix(filepath.Base(relativeDir), "dom")
}

var rules = []importRule{
	{
		// The old storage-shaped packages. They know nothing about the
		// terminal, the composition root, or the commands that drive them.
		//
		// internal/domain is deliberately absent from the forbid list: during
		// the migration a legacy shape may hold a domain value object (the
		// history index stores an operation.Operation). The reverse stays
		// forbidden -- the domain rule below allows a ring's own subtree only
		// -- so the arrow still points one way.
		name: "internal/data",
		dir:  "internal/data",
		forbid: []string{
			"internal/app",
			"internal/helper",
			"internal/gui",
			"internal/desktop",
			"cmd",
		},
	},
	{
		// The shell runner is infrastructure shared by every domain. It used to
		// be enough for one operation to import another; that is what this
		// stops.
		name:     "internal/shell",
		dir:      "internal/shell",
		required: true,
		forbid: []string{
			"internal/app",
			"internal/data",
			"internal/helper",
			"internal/gui",
			"internal/desktop",
			"cmd",
		},
	},
	{
		// Sentinel errors cross every layer, so this package must stay free of
		// dependencies; an import here is an import everywhere.
		name:       "internal/errors",
		dir:        "internal/errors",
		required:   true,
		noModule:   true,
		noExternal: true,
	},
	{
		// pkg is the public utility layer. It may use third-party code (the
		// TOML codec lives here) but it must never reach back into the
		// application.
		name:     "pkg",
		dir:      "pkg",
		required: true,
		forbid:   []string{"cmd"},
		noModule: false,
	},
	{
		// Ring 1. Nothing but the standard library: a domain package that
		// imports the filesystem, a codec or a neighbour is no longer a model.
		name:               "internal/domain",
		dir:                "internal/domain",
		noModuleOutsideOwn: true,
		noExternal:         true,
		skip:               legacyOperationPackage,
	},
	{
		// Ring 2. Use cases may reach the domain and their ports, never an
		// adapter or the driver.
		name: "internal/application",
		dir:  "internal/application",
		forbid: []string{
			"internal/adapters",
			"internal/app",
			"internal/data",
			"internal/gui",
			"internal/desktop",
			"cmd",
		},
	},
	{
		// Ring 3. Adapters implement ports; they do not drive the CLI.
		name: "internal/adapters",
		dir:  "internal/adapters",
		forbid: []string{
			"internal/app",
			"internal/data",
			"internal/gui",
			"internal/desktop",
			"cmd",
		},
	},
}

// TestRuleViolation pins the matcher. Without it the gate is only as good as
// its string comparison, and a prefix rule that silently stopped matching
// would report every package as clean.
func TestRuleViolation(t *testing.T) {
	shared := importRule{forbid: []string{"internal/app", "cmd"}}
	domain := importRule{dir: "internal/domain", noModuleOutsideOwn: true, noExternal: true}
	pure := importRule{noModule: true, noExternal: true}

	tests := []struct {
		name string
		rule importRule
		pkg  string
		imp  string
		want string
	}{
		{
			name: "forbidden package",
			rule: shared,
			imp:  modulePath + "/internal/app",
			want: "the dependency arrow points the other way",
		},
		{
			name: "forbidden subtree",
			rule: shared,
			imp:  modulePath + "/internal/app/save",
			want: "the dependency arrow points the other way",
		},
		{
			name: "a sibling with the forbidden name as a prefix is allowed",
			rule: shared,
			imp:  modulePath + "/internal/application",
			want: "",
		},
		{
			name: "the forbidden package as a prefix of a longer name is allowed",
			rule: importRule{forbid: []string{"cmd"}},
			imp:  modulePath + "/cmdutils",
			want: "",
		},
		{
			name: "another ring is allowed",
			rule: shared,
			imp:  modulePath + "/internal/shell",
			want: "",
		},
		{
			name: "the standard library is allowed",
			rule: shared,
			imp:  "fmt",
			want: "",
		},
		{
			name: "a domain may import itself",
			rule: domain,
			imp:  modulePath + "/internal/domain/catalog",
			want: "",
		},
		{
			name: "a domain may not import a neighbour",
			rule: domain,
			imp:  modulePath + "/internal/data/command",
			want: "this ring may only import from its own subtree",
		},
		{
			name: "a domain may not take a third-party dependency",
			rule: domain,
			imp:  "github.com/pelletier/go-toml",
			want: "this ring must depend on the standard library only",
		},
		{
			name: "a dependency-free package may not import the module",
			rule: pure,
			imp:  modulePath + "/internal/shell",
			want: "this package must not import the module",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.violation(tt.pkg, tt.imp); got != tt.want {
				t.Errorf("violation(%q) = %q, want %q", tt.imp, got, tt.want)
			}
		})
	}
}

// TestLayerDependencies is the gate itself.
func TestLayerDependencies(t *testing.T) {
	checkedAny := false

	for _, rule := range rules {
		packages := packagesUnder(t, rule.dir)

		applicable := 0

		for _, pkg := range packages {
			if rule.skip != nil && rule.skip(pkg.path) {
				continue
			}

			applicable++
		}

		if applicable == 0 {
			if rule.required {
				t.Errorf("%s: no package found under %q, so this rule checked nothing", rule.name, rule.dir)
			} else {
				t.Logf("%s: ring does not exist yet, rule inactive", rule.name)
			}

			continue
		}

		checkedAny = true
		t.Logf("%s: %d package(s)", rule.name, applicable)

		for _, pkg := range packages {
			if rule.skip != nil && rule.skip(pkg.path) {
				continue
			}

			for _, imported := range pkg.imports {
				if violation := rule.violation(pkg.path, imported); violation != "" {
					t.Errorf("%s (in %s) imports %s: %s", rule.name, pkg.path, imported, violation)
				}
			}
		}
	}

	if !checkedAny {
		t.Fatal("the architecture gate examined no package; it cannot fail and reports nothing")
	}
}

// violation explains which part of the rule an import breaks, or returns "".
func (r importRule) violation(pkgPath, imported string) string {
	if !strings.HasPrefix(imported, modulePath) {
		if r.noExternal && strings.Contains(firstPathElement(imported), ".") {
			return "this ring must depend on the standard library only"
		}

		return ""
	}

	rest := strings.TrimPrefix(strings.TrimPrefix(imported, modulePath), "/")

	if r.noModule {
		return "this package must not import the module"
	}

	if r.noModuleOutsideOwn {
		if r.matches(rest) {
			return ""
		}

		return "this ring may only import from its own subtree"
	}

	for _, forbidden := range r.forbid {
		if rest == forbidden || strings.HasPrefix(rest, forbidden+"/") {
			return "the dependency arrow points the other way"
		}
	}

	return ""
}

// matches reports whether a module-relative path is inside the rule's ring.
func (r importRule) matches(relativePath string) bool {
	return relativePath == r.dir || strings.HasPrefix(relativePath, r.dir+"/")
}

func firstPathElement(importPath string) string {
	first, _, _ := strings.Cut(importPath, "/")
	return first
}

// sourcePackage is one directory of Go source, with the imports of its
// non-test files.
type sourcePackage struct {
	path    string
	imports []string
}

// packagesUnder parses every non-test Go file below dir and returns one entry
// per directory that holds any.
//
// Test files are excluded on purpose: a test may build a fixture from an
// adapter while the code under test stays clean, and the rule is about the
// shipped binary's dependencies.
func packagesUnder(t *testing.T, dir string) []sourcePackage {
	t.Helper()

	fset := token.NewFileSet()
	byDir := map[string]map[string]bool{}

	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if name := entry.Name(); name == "testdata" || name == "vendor" {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		relative := filepath.ToSlash(filepath.Dir(path))

		if byDir[relative] == nil {
			byDir[relative] = map[string]bool{}
		}

		for _, spec := range file.Imports {
			imported := strings.Trim(spec.Path.Value, `"`)
			byDir[relative][imported] = true
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", dir, err)
	}

	packages := make([]sourcePackage, 0, len(byDir))

	for path, imports := range byDir {
		list := make([]string, 0, len(imports))
		for imported := range imports {
			list = append(list, imported)
		}

		packages = append(packages, sourcePackage{path: path, imports: list})
	}

	return packages
}
