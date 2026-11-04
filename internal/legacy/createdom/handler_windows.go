//go:build windows

package createdom

// CreateSystemVariable stores one variable in the user's environment.
func (h *Handler) CreateSystemVariable(newVar string, value string) (string, error) {
	return h.environment.SetVariable(newVar, value)
}

// CreateSystemPathVariable adds one directory to the user's Path value.
func (h *Handler) CreateSystemPathVariable(value string) (string, error) {
	return h.environment.AddPath(value)
}
