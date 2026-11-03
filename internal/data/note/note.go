package note

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/data/operations"
)

type Note struct {
	builder strings.Builder
}

func New() *Note {
	return &Note{
		builder: strings.Builder{},
	}
}

func (note *Note) AddCommandName(name string) *Note {
	fmt.Fprintf(&note.builder, ";command_name=%s", name)
	return note
}

func (note *Note) AddCommandAlias(aliases []string) *Note {
	fmt.Fprintf(&note.builder, ";command_alias=%v", aliases)
	return note
}

func (note *Note) AddCommandExec(exec string) *Note {
	fmt.Fprintf(&note.builder, ";command_exec=%s", exec)
	return note
}

func (note *Note) AddOldCommandName(name string) *Note {
	fmt.Fprintf(&note.builder, ";old_command_name=%s", name)
	return note
}

func (note *Note) AddOldCommandAlias(aliases []string) *Note {
	fmt.Fprintf(&note.builder, ";old_command_alias=%v", aliases)
	return note
}

func (note *Note) AddOldCommandExec(exec string) *Note {
	fmt.Fprintf(&note.builder, ";old_command_exec=%s", exec)
	return note
}

func (note *Note) AddCrowerExec(command string, args []string) *Note {
	var execBuilder strings.Builder
	// command is the operation and its flags without the program name, so the
	// recorded line reads "crower copy a b". A caller that included the name
	// produced "crower crower copy a b".
	fmt.Fprintf(&execBuilder, "crower %s", command)
	for _, arg := range args {
		fmt.Fprintf(&execBuilder, " %s", arg)
	}
	fmt.Fprintf(&note.builder, ";crower_exec=%s", execBuilder.String())
	return note
}

func (note *Note) AddMainOperation(op operations.MainOperationEnum) *Note {
	fmt.Fprintf(&note.builder, ";main_operation=%s", op)
	return note
}

func (note *Note) AddSubOperation(op operations.SubOperationEnum) *Note {
	fmt.Fprintf(&note.builder, ";sub_operation=%s", op)
	return note
}

func (note *Note) Build() string {

	return strings.Trim(note.builder.String(), ";")
}

// CommandName reads the command_name field out of a built note.
//
// A note is a ";"-separated list of key=value fields. History entries written
// before the name was stored on the entry itself carry it only here, so a run
// that loads an older history can still resolve it.
//
// Only a note for an operation that touches a managed command answers. A file
// or variable note can contain the same text inside a path, and resolving that
// would point a `--last` operation at a command the operation never touched.
func CommandName(built string) string {
	if !namesACommand(built) {
		return ""
	}

	for field := range strings.SplitSeq(built, ";") {
		if name, ok := strings.CutPrefix(field, "command_name="); ok {
			return name
		}
	}

	return ""
}

// namesACommand reports whether the note's sub-operation acts on a managed
// command, as opposed to a file, folder, variable or the history itself.
func namesACommand(built string) bool {
	for field := range strings.SplitSeq(built, ";") {
		switch field {
		case "sub_operation=command", "sub_operation=last", "sub_operation=script", "sub_operation=process":
			return true
		}
	}

	return false
}
