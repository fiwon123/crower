package payload

import "github.com/fiwon123/crower/internal/data/operation"

type Data struct {
	Op    operation.State
	Args  []string
	Name  string
	Alias []string
	Exec  string
}

func New(op operation.State, args []string, name string, alias []string, exec string) Data {
	return Data{
		Op:    op,
		Args:  args,
		Name:  name,
		Alias: alias,
		Exec:  exec,
	}
}
