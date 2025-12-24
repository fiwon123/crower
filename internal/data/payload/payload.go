package payload

import "github.com/fiwon123/crower/internal/data/operations"

type Data struct {
	Op    operations.Data
	Args  []string
	Name  string
	Alias []string
	Exec  string
}

func NewPayload(op operations.Data, args []string, name string, alias []string, exec string) Data {
	return Data{
		Op:    op,
		Args:  args,
		Name:  name,
		Alias: alias,
		Exec:  exec,
	}
}
