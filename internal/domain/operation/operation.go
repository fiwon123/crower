// Package operation names the operations crower records in its history.
//
// The values are a file format, not an internal label: each one is written as
// the Operation field of a historyhelper.json entry and as
// `main_operation=<value>` inside a history note. Renaming one makes the
// history of every existing user unreadable, so the strings are pinned by the
// tests next to this file.
package operation

// Operation is one of the operations the history records.
type Operation string

// SubOperation is what an operation acted on.
type SubOperation string

// The operations, in the order the help lists them.
const (
	Execute Operation = "execute"
	Create  Operation = "create"
	Delete  Operation = "delete"
	Update  Operation = "update"
	List    Operation = "list"
	Reset   Operation = "reset"
	Open    Operation = "open"
	Revert  Operation = "revert"
	Search  Operation = "search"
	Extract Operation = "extract"
	Copy    Operation = "copy"
	Move    Operation = "move"
	Restore Operation = "restore"
	Check   Operation = "check"
	Upgrade Operation = "upgrade"
)

// The subjects an operation can act on.
const (
	Command       SubOperation = "command"
	Script        SubOperation = "script"
	Process       SubOperation = "process"
	File          SubOperation = "file"
	Folder        SubOperation = "folder"
	FileAndFolder SubOperation = "file_and_folder"
	History       SubOperation = "history"
	SysPath       SubOperation = "syspath"
	System        SubOperation = "system"
	Last          SubOperation = "last"
	Browser       SubOperation = "browser"
)
