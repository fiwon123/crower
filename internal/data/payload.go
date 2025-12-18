package data

type Payload struct {
	Op      Operation
	Args    []string
	Command *Command
}
