package desktop

// HistoryText renders the history screen for another front end. The Ebiten
// window and the Fyne window show the same rows, so the rendering stays here
// where its test already lives.
func (m *Model) HistoryText() string {
	return m.historyText()
}
