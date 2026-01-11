package state

// all operations
type MainOperationEnum string
type SubOperationEnum string

// main
const (
	Execute MainOperationEnum = "execute"
	Create  MainOperationEnum = "create"
	Delete  MainOperationEnum = "delete"
	Update  MainOperationEnum = "update"
	List    MainOperationEnum = "list"
	Reset   MainOperationEnum = "reset"
	Open    MainOperationEnum = "open"
	Revert  MainOperationEnum = "revert"
	Search  MainOperationEnum = "search"
	Extract MainOperationEnum = "extract"
	Copy    MainOperationEnum = "copy"
	Move    MainOperationEnum = "move"
	Restore MainOperationEnum = "restore"
	Check   MainOperationEnum = "check"
)

// sub
const (
	Command       SubOperationEnum = "command"
	Process       SubOperationEnum = "process"
	File          SubOperationEnum = "file"
	Folder        SubOperationEnum = "folder"
	FileAndFolder SubOperationEnum = "file_and_folder"
	History       SubOperationEnum = "history"
	SysPath       SubOperationEnum = "syspath"
	System        SubOperationEnum = "system"
	Last          SubOperationEnum = "last"
)
