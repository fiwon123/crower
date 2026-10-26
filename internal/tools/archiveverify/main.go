// Command archiveverify checks the release archives in a directory.
//
// Every archive must contain the executables the documentation tells users to
// run — `crower` and `cr`, or `crower.exe` and `cr.exe` on Windows. The
// platform and architecture belong in the archive filename only; a binary
// called `cr_linux_amd64` is unreachable by the name the README uses, which is
// exactly how v0.4.0 shipped a release where `cr` did not exist.
//
// Usage:
//
//	go run ./internal/tools/archiveverify [dir]
//
// The directory defaults to build/dist. It exits 1 when any archive is missing
// a name, carries an arch-suffixed binary, or when the directory holds no
// archive at all — an empty check must not read as a pass.
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	appName   = "crower"
	shortName = "cr"
)

func main() {
	dir := filepath.Join("build", "dist")
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	archives, err := findArchives(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "archiveverify: %v\n", err)
		os.Exit(1)
	}

	failed := 0
	for _, path := range archives {
		names, err := listArchive(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", path, err)
			failed++
			continue
		}
		problems := checkNames(names)
		if len(problems) > 0 {
			fmt.Fprintf(os.Stderr, "FAIL %s\n", path)
			for _, p := range problems {
				fmt.Fprintf(os.Stderr, "  - %s\n", p)
			}
			failed++
			continue
		}
		fmt.Printf("ok   %s (crower, cr)\n", path)
	}

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "\narchiveverify: %d of %d archive(s) are wrong\n", failed, len(archives))
		os.Exit(1)
	}
	fmt.Printf("\narchiveverify: %d archive(s) carry crower and cr\n", len(archives))
}

// findArchives lists the release archives in dir. No archive is an error: a
// check that inspects nothing would otherwise report success.
func findArchives(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	var archives []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip") {
			archives = append(archives, filepath.Join(dir, name))
		}
	}
	if len(archives) == 0 {
		return nil, fmt.Errorf("no .tar.gz or .zip archives in %s", dir)
	}
	sort.Strings(archives)
	return archives, nil
}

// checkNames reports what is wrong with the entries of one archive.
func checkNames(entries []string) []string {
	have := make(map[string]bool, len(entries))
	var suffixed []string
	for _, entry := range entries {
		base := filepath.Base(entry)
		have[base] = true
		if isArchSuffixed(base) {
			suffixed = append(suffixed, base)
		}
	}

	var problems []string
	if !have[appName] && !have[appName+".exe"] {
		problems = append(problems, fmt.Sprintf(
			"missing %q (or %q); archive holds %v", appName, appName+".exe", entries,
		))
	}
	if !have[shortName] && !have[shortName+".exe"] {
		problems = append(problems, fmt.Sprintf(
			"missing %q (or %q); archive holds %v", shortName, shortName+".exe", entries,
		))
	}
	for _, name := range suffixed {
		problems = append(problems, fmt.Sprintf(
			"arch-suffixed binary %q: package the plain name, the archive filename already carries the platform", name,
		))
	}
	return problems
}

// isArchSuffixed reports whether base looks like crower_linux_amd64 or
// cr_windows_arm64.exe — a binary named after its platform instead of its
// command.
func isArchSuffixed(base string) bool {
	stem := strings.TrimSuffix(strings.TrimSuffix(base, ".exe"), ".exe")
	if !strings.HasPrefix(stem, appName+"_") && !strings.HasPrefix(stem, shortName+"_") {
		return false
	}
	for _, targetOS := range []string{"linux", "windows", "darwin", "freebsd"} {
		if strings.Contains(stem, "_"+targetOS) {
			return true
		}
	}
	return false
}

// listArchive returns the entry names of a .tar.gz or .zip archive.
func listArchive(path string) ([]string, error) {
	switch {
	case strings.HasSuffix(path, ".zip"):
		return listZip(path)
	case strings.HasSuffix(path, ".tar.gz"):
		return listTarGz(path)
	default:
		return nil, errors.New("unsupported archive type")
	}
}

func listZip(path string) ([]string, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	names := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		names = append(names, file.Name)
	}
	return names, nil
}

func listTarGz(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var names []string
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		names = append(names, header.Name)
	}
	return names, nil
}
