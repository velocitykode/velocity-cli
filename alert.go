package prism

// Note prints a message in a rounded box with primary-colored border.
func Note(message string) {
	fprintln(noteBoxStyle.Render(message))
}

// Alert prints a message in a rounded box with error-colored border.
func Alert(message string) {
	fprintln(alertBoxStyle.Render(message))
}
