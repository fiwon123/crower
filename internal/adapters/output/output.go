// Package output is the technology side of a run's I/O: the three streams the
// program reads and writes, the result format, and the encoding of a
// structured result.
//
// It knows nothing about commands. A use case describes a result; this package
// decides what bytes reach the terminal, and `pkg/logger` decides which stream
// and which log they belong to.
package output
