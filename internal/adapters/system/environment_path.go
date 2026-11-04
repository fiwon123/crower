package system

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/pkg/utils"
)

// addPathToFile adds the value to the PATH export of the file. It reports
// whether the file was changed: a file that already has the value in its
// PATH is left alone.
func addPathToFile(filePath string, value string) (bool, error) {
	lineSlice := utils.GetFileLineSlice(filePath)

	for i, line := range lineSlice {
		if !isPathExport(line) {
			continue
		}

		if pathHasValue(line, value) {
			return false, nil
		}

		lineSlice[i] = addValueToPathExport(line, value)

		if err := utils.WriteFile(lineSlice, filePath); err != nil {
			return false, err
		}

		return true, nil
	}

	lineSlice = append(lineSlice, fmt.Sprintf("export PATH=\"$PATH:%s\"", value))

	if err := utils.WriteFile(lineSlice, filePath); err != nil {
		return false, err
	}

	return true, nil
}

// isPathExport reports whether the line is a PATH export. A line that starts
// with "export PATH" but is not one — "export PATHDIR=…" — is not.
func isPathExport(line string) bool {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "export PATH=") {
		return true
	}

	return trimmed == "export PATH" || strings.HasPrefix(trimmed, "export PATH ")
}

// pathHasValue reports whether the PATH export already has the value as one
// of its elements. The elements are compared whole: "/opt/tools" is not in a
// PATH that holds "/opt/tools-extra".
func pathHasValue(line, value string) bool {
	for _, element := range pathElements(line) {
		if element == value {
			return true
		}
	}

	return false
}

// pathElements returns the elements of the PATH export of the line.
func pathElements(line string) []string {
	trimmed := strings.TrimSpace(line)

	value := ""
	if index := strings.Index(trimmed, "="); index >= 0 {
		value = trimmed[index+1:]
	}

	value = strings.Trim(value, "\"")

	return strings.Split(value, ":")
}

// addValueToPathExport returns the PATH export with the value added to it,
// keeping the form the line already uses.
func addValueToPathExport(line, value string) string {
	trimmed := strings.TrimSpace(line)

	// A line that is only "export PATH" carries no value to extend.
	if !strings.Contains(trimmed, "=") {
		return fmt.Sprintf("export PATH=\"$PATH:%s\"", value)
	}

	start := strings.Index(trimmed, "\"")
	end := strings.LastIndex(trimmed, "\"")

	if start >= 0 && end > start {
		return trimmed[:end] + ":" + value + trimmed[end:]
	}

	return trimmed + ":" + value
}
