package desktop

// Groups stores the icon each group shows.
//
// It is declared here, next to Tools and Executor, so this package depends on
// neither the storage adapter nor the command layer; cmd supplies one and
// internal/fyneui stays build-tagged.
type Groups interface {
	// Icons returns the icon of every group that has one.
	Icons() (map[string]string, error)

	// SetIcon stores or, with an empty icon, removes one group's icon and
	// returns the whole map.
	SetIcon(group, icon string) (map[string]string, error)
}
