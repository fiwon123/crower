//go:build linux

package createdom

// CreateSystemVariable adds one variable to the file a login shell reads.
func (h *Handler) CreateSystemVariable(newVar string, value string) (string, error) {
	return h.environment.SetVariable(newVar, value)
}

// CreateSystemPathVariable adds one directory to the PATH of the user.
func (h *Handler) CreateSystemPathVariable(value string) (string, error) {
	return h.environment.AddPath(value)
}
