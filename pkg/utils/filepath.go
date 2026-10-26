package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// CreateFileIfNotExists creates filePath when it is missing and leaves an
// existing file untouched.
//
// It returns an error rather than calling log.Fatal: the caller knows whether a
// missing file is fatal, and a util that exits the process takes that decision
// away from it.
func CreateFileIfNotExists(filePath string) error {
	file, err := os.OpenFile(filePath,
		os.O_CREATE|os.O_RDWR,
		0644)
	if err != nil {
		return fmt.Errorf("create file %s: %w", filePath, err)
	}

	return file.Close()
}

func CreateFolderIfNotExists(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create directory %s: %w", path, err)
	}

	return nil
}

func DeleteFile(path string) error {
	err := os.RemoveAll(path)
	return err
}

func FilePathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func FolderPathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

// Consider path if has extension
func IsValidFilePath(path string) bool {
	if path == "" {
		return false
	}

	if path == "." {
		return false
	}

	if !strings.Contains(path, "/") && !strings.Contains(path, "\\") {
		return false
	}

	base := filepath.Base(path)
	splitted := strings.Split(base, ".")

	return len(splitted) == 2
}

func IsValidFolderPath(path string) bool {
	if path == "" {
		return false
	}

	if path == "." {
		return false
	}

	if !strings.Contains(path, "/") && !strings.Contains(path, "\\") {
		return false
	}

	return true
}

func GetFileLineSlice(filePath string) []string {

	lines, err := os.ReadFile(filePath)
	if err != nil {
		lines = []byte{}
	}

	lineSlice := strings.Split(string(lines), "\n")

	return lineSlice
}

func WriteFile(lineSlice []string, filePath string) error {

	err := os.WriteFile(filePath, []byte(strings.Join(lineSlice, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("Error writing .profile: %v", err)
	}

	return nil
}

func LineExists(filePath, line string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == line {
			return true
		}
	}
	return false
}

func SplitPath(path string) []string {
	return filepath.SplitList(path)
}

func CheckNewVarValuePath(value string, from string) bool {
	splitted := SplitPath(from)
	ok := true
	if slices.Contains(splitted, value) {
		return false
	}

	return ok
}
