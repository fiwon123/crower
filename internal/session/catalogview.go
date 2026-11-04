package session

import (
	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// Catalog returns the managed commands as the domain aggregate.
//
// It is the bridge the new use cases use while the map-based core is still the
// storage: the aggregate is built on every call, so a caller that holds one
// does not see a later mutation through it.
func (s *Session) Catalog() catalog.Catalog {
	return catalogOf(s.allCommandsByName, s.orderKeys)
}

// ReplaceCatalog writes the aggregate back into the maps the operations read.
//
// Both maps are rebuilt, not just the name map. Clearing one and leaving the
// other is what `reset` used to do: the alias index kept resolving commands
// that no longer existed, which a command that continues to run after the
// reset would see.
func (s *Session) ReplaceCatalog(commands catalog.Catalog) {
	s.allCommandsByName = command.NewMapData()
	s.orderKeys = commands.Names()

	for _, cmd := range commands.Commands() {
		s.allCommandsByName.Add(cmd.Name, toData(cmd))
	}

	// The alias index is rebuilt with the loader's rule -- names sorted, first
	// claim wins -- rather than by iteration order, so a config that already
	// holds a collision resolves the same way before and after a change.
	s.allCommandsByAlias = aliasMapFrom(commands)
}
