package app

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
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

// withDefaults fills any stream that was left nil, so a caller that only cares
// about one stream does not have to construct the other two.
func (s Streams) withDefaults() Streams {
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

// Options are the settings a command resolves from its global flags and hands
// to InitApp.
type Options struct {
	Streams Streams

	// NoInput makes every prompt fail instead of reading stdin, so a
	// non-interactive run can never block waiting for a person.
	NoInput bool
}

// withDefaults fills in the terminal a run would use by default.
func (o Options) withDefaults() Options {
	o.Streams = o.Streams.withDefaults()
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
