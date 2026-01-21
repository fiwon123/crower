package notes

import (
	"fmt"
	"strings"

	operationsdata "github.com/fiwon123/crower/internal/data/operations"
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
	fmt.Fprintf(&execBuilder, "crower %s", command)
	for _, arg := range args {
		fmt.Fprintf(&execBuilder, " %s", arg)
	}
	execBuilder.WriteString(" ")

	fmt.Fprintf(&note.builder, ";crower_exec=%s", execBuilder.String())
	return note
}

func (note *Note) AddMainOperation(op operationsdata.MainOperationEnum) *Note {
	fmt.Fprintf(&note.builder, ";main_operation=%s", op)
	return note
}

func (note *Note) AddSubOperation(op operationsdata.SubOperationEnum) *Note {
	fmt.Fprintf(&note.builder, ";sub_operation=%s", op)
	return note
}

func (note *Note) Build() string {

	return strings.Trim(note.builder.String(), ";")
}
