package data

type Operation string

const (
	ExecuteOp  Operation = "execute"
	AddOp      Operation = "add"
	ProcessOp  Operation = "process"
	AddProcess Operation = "addProcess"
	DeleteOp   Operation = "delete"
	UpdateOp   Operation = "update"
	ListOp     Operation = "list"
	ResetOp    Operation = "reset"
	OpenOp     Operation = "open"
	HistoryOp  Operation = "history"
	RevertOp   Operation = "revert"
)
