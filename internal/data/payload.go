package data

type Payload struct {
	Op    Operation
	Args  []string
	Name  string
	Alias []string
	Exec  string
}
