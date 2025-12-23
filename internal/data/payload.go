package data

type Payload struct {
	Op    Operation
	Args  []string
	Name  string
	Alias []string
	Exec  string
}

func NewPayload(op Operation, args []string, name string, alias []string, exec string) Payload {
	return Payload{
		Op:    op,
		Args:  args,
		Name:  name,
		Alias: alias,
		Exec:  exec,
	}
}
