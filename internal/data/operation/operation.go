package operation

type State string

const (
	Execute    State = "execute"
	Create     State = "create"
	Process    State = "process"
	AddProcess State = "addProcess"
	Delete     State = "delete"
	Update     State = "update"
	List       State = "list"
	Reset      State = "reset"
	Open       State = "open"
	History    State = "history"
	Revert     State = "revert"
)
