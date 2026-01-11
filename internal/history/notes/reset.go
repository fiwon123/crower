package notes

// Create a new reset note
func GenerateResetNote() string {
	noteBuilder := New()
	noteBuilder.AddCrowerExec("reset", nil)

	return noteBuilder.Build()
}
