package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Format is how an operation with a structured result renders it.
type Format string

const (
	// FormatText writes the rendering the operation built for a person. It is
	// the default, and it is the only format the CLI had before --json.
	FormatText Format = "text"

	// FormatJSON writes one JSON object for a script to read.
	FormatJSON Format = "json"
)

// Result writes one operation result in the configured format.
//
// Text mode writes text verbatim: the operation already knows what a table, a
// listing or a search should look like, and rewriting it here would change the
// bytes a script currently reads. JSON mode ignores text and encodes payload
// instead, so a result has exactly one owner per format and both go through
// Logger.Out -- the result stream and the log file see the same bytes.
//
// Raw command output (execute, extract, open) is deliberately not routed
// through here: those bytes belong to the command that produced them, and
// wrapping them in an envelope would break `crower execute ... | grep`.
func (a *Config) Result(text string, payload any) error {
	if a.Format != FormatJSON {
		a.Logger.Out(text)
		return nil
	}

	var buffer bytes.Buffer

	encoder := json.NewEncoder(&buffer)
	// A path is data, not HTML. Without this, `a&b` is encoded as
	// `a\u0026b`: valid JSON, and a value a script cannot match against the
	// path it passed in.
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(payload); err != nil {
		return fmt.Errorf("encode json result: %w", err)
	}

	// Encoder terminates the object; Out adds the newline a result is expected
	// to end with, so one result stays exactly one line.
	a.Logger.Out(strings.TrimRight(buffer.String(), "\n"))

	return nil
}
