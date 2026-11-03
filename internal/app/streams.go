package app

import (
	"bufio"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/output"
	"github.com/fiwon123/crower/internal/errors"
)

// Streams is the I/O surface of a run. It is the output adapter's type: the
// application names the interface it needs, and the technology provides it.
type Streams = output.Streams

// DefaultStreams is the terminal a binary runs against.
func DefaultStreams() Streams {
	return output.DefaultStreams()
}

// Options are the settings a command resolves from its global flags and hands
// to InitApp.
type Options struct {
	Streams Streams

	// NoInput makes every prompt fail instead of reading stdin, so a
	// non-interactive run can never block waiting for a person.
	NoInput bool

	// Format is how an operation with a structured result renders it.
	Format Format
}

// withDefaults fills in the terminal a run would use by default.
func (o Options) withDefaults() Options {
	o.Streams = o.Streams.WithDefaults()

	if o.Format == "" {
		o.Format = FormatText
	}

	return o
}

// ReadLine reads one line from the input stream.
//
// A closed stream is an error, not an empty answer: a prompt that received
// nothing must fail, and the old reader ignored the error and asked again
// forever, so `crower create </dev/null` hung and wrote hundreds of thousands
// of lines.
func (a *Config) ReadLine() (string, error) {
	line, err := a.in.ReadString('\n')
	if err != nil && line == "" {
		return "", errors.NoInputAvailable()
	}

	return strings.TrimRight(line, "\r\n"), nil
}

// Reader is the buffered input stream, so several prompts in one command share
// one buffer. Wrapping os.Stdin per prompt used to discard whatever the last
// prompt had already buffered.
func (a *Config) Reader() *bufio.Reader {
	return a.in
}
