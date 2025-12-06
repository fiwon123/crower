package data

type CommandOperation string

const (
	Execute CommandOperation = "execute"
	Create  CommandOperation = "create"
	Delete  CommandOperation = "delete"
	Update  CommandOperation = "update"
	List    CommandOperation = "list"
	Reset   CommandOperation = "reset"
)
