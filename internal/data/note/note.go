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
	fmt.Fprintf(&note.builder, ";command_name=%s", escapeValue(name))
	return note
}

func (note *Note) AddCommandAlias(aliases []string) *Note {
	fmt.Fprintf(&note.builder, ";command_alias=%s", escapeValue(fmt.Sprintf("%v", aliases)))
	return note
}

func (note *Note) AddCommandExec(exec string) *Note {
	fmt.Fprintf(&note.builder, ";command_exec=%s", escapeValue(exec))
	return note
}

func (note *Note) AddOldCommandName(name string) *Note {
	fmt.Fprintf(&note.builder, ";old_command_name=%s", escapeValue(name))
	return note
}

func (note *Note) AddOldCommandAlias(aliases []string) *Note {
	fmt.Fprintf(&note.builder, ";old_command_alias=%s", escapeValue(fmt.Sprintf("%v", aliases)))
	return note
}

func (note *Note) AddOldCommandExec(exec string) *Note {
	fmt.Fprintf(&note.builder, ";old_command_exec=%s", escapeValue(exec))
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
	fmt.Fprintf(&note.builder, ";crower_exec=%s", escapeValue(execBuilder.String()))
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

// escapeValue encodes the two characters a note value cannot carry raw: the ";"
// that separates fields, and the "%" that introduces an escape. Every value a
// caller supplies goes through here, so an exec like "make; echo done" survives
// the round trip instead of being cut at the semicolon.
//
// A note written before this encoding still reads: a raw value produces neither
// "%25" nor "%3B" by accident, except for a literal "%3B" (or "%25") that a user
// typed on purpose. Those rare sequences decode to ";" (or "%"); the alternative
// -- a new note format -- would need a migration for every stored history.
func escapeValue(value string) string {
	value = strings.ReplaceAll(value, "%", "%25")
	return strings.ReplaceAll(value, ";", "%3B")
}

// unescapeValue reverses escapeValue. Decoding "%3B" before "%25" is what keeps
// a literal "%3B" (written as "%253B") from turning into ";" on the way back.
func unescapeValue(value string) string {
	value = strings.ReplaceAll(value, "%3B", ";")
	return strings.ReplaceAll(value, "%25", "%")
}

// Display renders a built note for a human, decoding each value so the escaped
// ";" and "%" read as they were written. It is for output only: read a field
// with CommandName or CommandExec, never by splitting Display's result.
func Display(built string) string {
	var out strings.Builder
	first := true

	for field := range strings.SplitSeq(built, ";") {
		if !first {
			out.WriteString(";")
		}
		first = false

		key, value, ok := strings.Cut(field, "=")
		if !ok {
			out.WriteString(field)
			continue
		}

		out.WriteString(key)
		out.WriteString("=")
		out.WriteString(unescapeValue(value))
	}

	return out.String()
}

// CommandExec reads the command_exec field out of a built note.
//
// The execute note that `create --execute` reuses records the command's exec
// under command_exec. The reader used to look for a bare "exec" key, which no
// note builder has ever written, so the reuse could only ever fail.
//
// Like CommandName, it only answers for a note whose sub-operation acts on a
// managed command: a file or folder note can carry the same text inside a
// quoted path, and reading that would hand back a line the operation never ran.
func CommandExec(built string) string {
	if !namesACommand(built) {
		return ""
	}

	for field := range strings.SplitSeq(built, ";") {
		if exec, ok := strings.CutPrefix(field, "command_exec="); ok {
			return unescapeValue(exec)
		}
	}

	return ""
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
			return unescapeValue(name)
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
