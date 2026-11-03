package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// EncodeJSON encodes one result as a single line, ready for the result stream.
//
// A path is data, not HTML: without SetEscapeHTML(false), `a&b` is encoded as
// `a\u0026b`, which is valid JSON and a value a script cannot match against the
// path it passed in.
//
// The encoder terminates its object with a newline; that is trimmed here
// because the writer adds the one newline a result is expected to end with, and
// a result has to stay exactly one line for `| jq` and `| wc -l`.
func EncodeJSON(payload any) (string, error) {
	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(payload); err != nil {
		return "", fmt.Errorf("encode json result: %w", err)
	}

	return strings.TrimRight(buffer.String(), "\n"), nil
}
