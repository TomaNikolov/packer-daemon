package printer

import (
	"fmt"
)

// ConsolePrinter ...
type ConsolePrinter struct {
}

// NewConsolePrinter ...
func NewConsolePrinter() ConsolePrinter {
	return ConsolePrinter{}
}

// Print ...
func (consolePrinter ConsolePrinter) Print(message string) {
	fmt.Println(message)
}
