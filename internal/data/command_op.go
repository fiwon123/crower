package data

type Operation string

const (
	Execute Operation = "execute"
	Add     Operation = "add"
	Delete  Operation = "delete"
	Update  Operation = "update"
	List    Operation = "list"
	Reset   Operation = "reset"
	Open    Operation = "open"
	Process Operation = "process"
)
