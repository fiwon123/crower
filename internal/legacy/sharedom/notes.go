package sharedom

import "github.com/fiwon123/crower/internal/data/note"

func getCreateNote() string {
	nt := &note.Note{}
	nt.AddMainOperation("share")
	nt.AddSubOperation("create")
	nt.AddCommandName("")
	nt.AddCommandAlias([]string{})
	nt.AddCommandExec("")
	return nt.Build()
}
