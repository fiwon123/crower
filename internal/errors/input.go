package errors

import "fmt"

// InputRequired reports a prompt that would be shown but could not be, because
// --no-input was given. It is an error rather than an empty answer so a script
// never proceeds on a value a person was supposed to choose.
func InputRequired(ask string) error {
	return fmt.Errorf("input required for %q but --no-input is set", ask)
}

// NoInputAvailable reports that the input stream ended before an answer
// arrived. It is an error, not an empty answer, so a prompt cannot loop: the
// old implementation ignored the read error and re-prompted forever.
func NoInputAvailable() error {
	return fmt.Errorf("no input available: standard input is closed")
}

// InvalidRows reports a input that is not a row of the printed table.
func InvalidRows() error {
	return fmt.Errorf("invalid row")
}

// InvalidInput reports a confirmation that is neither yes nor no.
func InvalidInput() error {
	return fmt.Errorf("invalid input")
}

// EmptyInput reports a value that was required but left empty.
func EmptyInput() error {
	return fmt.Errorf("input is empty")
}

// OnlyLettersAndNumbers reports an alias with a character an alias cannot
// carry.
func OnlyLettersAndNumbers() error {
	return fmt.Errorf("only letters and numbers are allowed")
}
