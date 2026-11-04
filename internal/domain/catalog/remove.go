package catalog

// Remove deletes the command a key resolves to.
//
// The key resolves exactly as it does for Update -- a name first, then an
// alias -- so a delete cannot reach a different command than an update would.
// The receiver is unchanged and the returned catalog is a new value.
//
// The alias index is rebuilt from the commands that remain. A command that also
// declared an alias of the deleted one therefore takes that alias over at once,
// which is the answer a reload of the written file gives (the loader rebuilds
// the index from the file). Removing the only declarer leaves the alias
// unclaimed, because nothing is left to claim it.
func (c Catalog) Remove(key string) (Catalog, Command, bool) {
	command, found := c.Resolve(key)
	if !found {
		return c, Command{}, false
	}

	remaining := make([]Command, 0, len(c.order)-1)

	for _, name := range c.order {
		if name == command.Name {
			continue
		}

		remaining = append(remaining, c.byName[name])
	}

	// c.order still names the deleted command: New ignores a name the catalog
	// does not hold, so the surviving commands keep their display order.
	return New(c.order, remaining), command, true
}
