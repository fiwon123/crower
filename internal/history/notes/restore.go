package notes

// Create a new restore note
func GenerateRestoreNote(msg string) string {
	noteBuilder := New()
	noteBuilder.AddCrowerExec("restore", nil)

	return noteBuilder.Build()

}
