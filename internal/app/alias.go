package app

import (
	"maps"
	"slices"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// catalogOf adapts the persisted command map to the domain aggregate.
//
// The map is what the configuration file holds and what the operations still
// mutate; the aggregate is where the rules about it live. This function is the
// one place that knows both shapes, so the rules have a single implementation
// and the loader keeps a single call site.
func catalogOf(commandMap command.MapData, order []string) catalog.Catalog {
	commands := make([]catalog.Command, 0, len(commandMap))

	for _, name := range slices.Sorted(maps.Keys(commandMap)) {
		stored := commandMap[name]
		commands = append(commands, catalog.Command{
			Name:    stored.Name,
			Aliases: stored.AllAlias,
			Exec:    stored.Exec,
			Steps:   stored.Steps,
		})
	}

	return catalog.New(order, commands)
}

// aliasMapFrom renders the aggregate's alias index as the map the operations
// look commands up in.
func aliasMapFrom(c catalog.Catalog) command.MapData {
	aliases := command.NewMapData()

	for alias, name := range c.AliasIndex() {
		owned, ok := c.Get(name)
		if !ok {
			continue
		}

		aliases.Add(alias, toData(owned))
	}

	return aliases
}

// getAliasMap builds the alias index in a stable order: commands sorted by
// name, first claim wins.
func getAliasMap(commandMap command.MapData) command.MapData {
	return aliasMapFrom(catalogOf(commandMap, nil))
}

// commandCollisions reports every key two commands claim, and every alias that
// is also another command's name.
//
// It is kept as a function here because the loader and the tests call it, but
// the report itself belongs to the catalog: the aggregate is what knows which
// claim wins.
func commandCollisions(commandMap command.MapData) []string {
	return catalogOf(commandMap, nil).Warnings()
}

// toData converts an aggregate command back to the persisted shape.
func toData(c catalog.Command) *command.Data {
	return &command.Data{
		Name:     c.Name,
		AllAlias: c.Aliases,
		Exec:     c.Exec,
		Steps:    c.Steps,
	}
}
