package desktop

// Groups stores the groups the window offers and the icon each shows.
//
// It is declared here, next to Tools and Executor, so this package depends on
// neither the storage adapter nor the command layer; cmd supplies one and
// internal/fyneui stays build-tagged.
type Groups interface {
	// Icons returns every stored group with its icon. A group with no icon is
	// present with an empty value: the group is the record, the icon decoration.
	Icons() (map[string]string, error)

	// SetIcon stores one group and its icon, creating the group when it is new,
	// and returns the whole map. An empty icon keeps the group.
	SetIcon(group, icon string) (map[string]string, error)

	// Remove deletes one group and returns the whole map.
	Remove(group string) (map[string]string, error)
}
