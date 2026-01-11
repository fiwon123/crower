package notes

// Create a new check note
func GenerateCheckNote() string {
	noteBuilder := New()
	noteBuilder.AddCrowerExec("crower --check", nil)

	return noteBuilder.Build()
}
