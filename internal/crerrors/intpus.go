package crerrors

import "fmt"

func InvalidRows() error {
	return fmt.Errorf("Invalid row")
}

func InvalidInput() error {
	return fmt.Errorf("Invalid input")
}

func EmptyInput() error {
	return fmt.Errorf("input is empty")
}

func OnlyLettersAndNumbers() error {
	return fmt.Errorf("Only numbers and letters")
}
