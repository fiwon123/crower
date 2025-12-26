package operation

type State string

const (
	Execute       State = "execute"
	ExecuteLast   State = "executeLast"
	ExecuteCreate State = "executeCreate"
	ExecuteUpdate State = "executeUpdate"
	Create        State = "create"
	CreateProcess State = "createProcess"
	Delete        State = "delete"
	DeleteCreate  State = "deleteCreate"
	DeleteUpdate  State = "deleteUpdate"
	DeleteExecute State = "deleteExecute"
	Update        State = "update"
	UpdateLast    State = "updateLast"
	UpdateCreate  State = "updateCreate"
	UpdateExecute State = "updateExecute"
	List          State = "list"
	ListProcess   State = "listProcess"
	ListHistory   State = "listHistory"
	Reset         State = "reset"
	Open          State = "open"
	OpenFolder    State = "openFolder"
	Revert        State = "revert"
)
