package notes

// Create a new upgrade note
func GenerateUpgradeNote() string {
	noteBuilder := New()
	noteBuilder.AddCrowerExec("--upgrade", nil)

	return noteBuilder.Build()
}
