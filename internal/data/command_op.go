package data

type CommandOperation string

const (
	Execute CommandOperation = "execute"
	Add     CommandOperation = "add"
	Delete  CommandOperation = "delete"
	Update  CommandOperation = "update"
	List    CommandOperation = "list"
	Reset   CommandOperation = "reset"
)
