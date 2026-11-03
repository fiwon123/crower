package output

import (
	"io"
	"os"
)

// Streams is the I/O surface of a crower run.
//
// Results go to Out; everything else -- prompts, status lines, warnings,
// errors -- is a diagnostic and goes to Err. Keeping them apart is what makes
// `crower list | wc -l` count rows instead of counting log noise, and it means
// a command's output is never interleaved with a status line about it.
//
// The fields are interfaces so a test can drive the whole command layer with
// buffers instead of a terminal.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// DefaultStreams is the terminal a binary runs against.
func DefaultStreams() Streams {
	return Streams{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}
}

// WithDefaults fills any stream that was left nil, so a caller that only cares
// about one stream does not have to construct the other two.
func (s Streams) WithDefaults() Streams {
	if s.In == nil {
		s.In = os.Stdin
	}

	if s.Out == nil {
		s.Out = os.Stdout
	}

	if s.Err == nil {
		s.Err = os.Stderr
	}

	return s
}

// Format is how an operation with a structured result renders it.
type Format string

const (
	// FormatText writes the rendering the operation built for a person. It is
	// the default, and it is the only format the CLI had before --json.
	FormatText Format = "text"

	// FormatJSON writes one JSON object for a script to read.
	FormatJSON Format = "json"
)
