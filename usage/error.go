package usage

import (
	"fmt"
)

// Error indicates usage information should be printed and exit with status code 2
type Error error

func GenericError() error {
	return Error(fmt.Errorf(""))
}

// Errorf creates a usage error with the given message
func Errorf(message string, formatArgs ...interface{}) error {
	return Error(fmt.Errorf(message, formatArgs...))
}
