package system

// Environment edits the environment variables the user's shell reads at start.
//
// Which file holds them belongs to the platform: Linux keeps a variable in
// ~/.bashrc and PATH in ~/.profile, macOS keeps both in ~/.zshrc, and Windows
// writes the user's Environment registry key. The operation says what to
// remove; this package knows where the host keeps it.
type Environment struct{}

// NewEnvironment builds an Environment that edits the host's files.
func NewEnvironment() *Environment {
	return &Environment{}
}
