package operations

type Data string

const (
	Execute    Data = "execute"
	Add        Data = "add"
	Process    Data = "process"
	AddProcess Data = "addProcess"
	Delete     Data = "delete"
	Update     Data = "update"
	List       Data = "list"
	Reset      Data = "reset"
	Open       Data = "open"
	History    Data = "history"
	Revert     Data = "revert"
)
