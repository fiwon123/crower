package desktop

import (
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// State is what a window shows after a setup switch: the commands and the
// history that belong to the setup now in use.
type State struct {
	Catalog catalog.Catalog
	History timeline.Log
}

// Setups lists and switches the named setups.
//
// It is declared here, next to Tools and Groups, so this package depends on
// neither the session nor the command layer; cmd supplies one and the UI
// package stays build-tagged.
type Setups interface {
	// Current is the setup in use, "default" when none was named.
	Current() string

	// Names lists every available setup, "default" first.
	Names() []string

	// Switch makes name the default for future runs and reloads the session,
	// returning the state the window should show.
	Switch(name string) (State, error)

	// Commands counts the commands a setup holds, so the selector can show what
	// switching would bring. It must not create a setup that does not exist.
	Commands(name string) (int, error)

	// Create makes an empty named setup, the way `crower setup new` does.
	Create(name string) error

	// Delete removes a named setup and its history. The default setup and the
	// active one are refused -- the same rules `crower setup delete` applies.
	Delete(name string) error
}
