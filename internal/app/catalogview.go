package app

import (
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// Catalog returns the managed commands as the domain aggregate.
//
// It is the bridge the new use cases use while the map-based core is still the
// storage: the aggregate is built on every call, so a caller that holds one
// does not see a later mutation through it.
func (app *Config) Catalog() catalog.Catalog {
	return catalogOf(app.AllCommandsByName, app.OrderKeys)
}

// ReplaceCatalog writes the aggregate back into the maps the operations read.
//
// Both maps are rebuilt, not just the name map. Clearing one and leaving the
// other is what `reset` used to do: the alias index kept resolving commands
// that no longer existed, which a command that continues to run after the
// reset would see.
func (app *Config) ReplaceCatalog(commands catalog.Catalog) {
	app.AllCommandsByName = command.NewMapData()
	app.AllCommandsByAlias = command.NewMapData()
	app.OrderKeys = commands.Names()

	for _, cmd := range commands.Commands() {
		stored := toData(cmd)
		app.AllCommandsByName.Add(cmd.Name, stored)

		for _, alias := range cmd.Aliases {
			app.AllCommandsByAlias.Add(alias, stored)
		}
	}
}
