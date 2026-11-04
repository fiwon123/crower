package system

import "github.com/fiwon123/crower/pkg/utils"

// Paths answers questions about the local filesystem.
type Paths struct{}

// NewPaths builds a Paths that reads the host.
func NewPaths() *Paths {
	return &Paths{}
}

// IsFile reports whether path exists and is a file.
//
// It is the question the transfer asks to tell a file from a folder, and it is
// deliberately not "exists": a folder exists too, and the legacy command sent
// anything that was not a file down the folder branch.
func (p *Paths) IsFile(path string) bool {
	return utils.FilePathExists(path)
}
