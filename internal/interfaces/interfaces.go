package interfaces

// ListHandler renders the command table the prompt screens show.
//
// It returns text instead of printing it, so the caller decides which stream
// receives the table: a result for `list`, a diagnostic for a prompt.
type ListHandler interface {
	CommandsTable() string
}
