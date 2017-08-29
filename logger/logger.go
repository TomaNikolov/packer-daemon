package logger

import (
	"fmt"

	"github.com/tomanikolov/packer-daemon/types"
)

// Logger ...
type Logger struct {
	Printers []types.Printer
}

// NewLogger ...
func NewLogger(printers []types.Printer) Logger {
	return Logger{
		Printers: printers,
	}
}

// Log ..
func (logger *Logger) Log(message string) {
	for _, printer := range logger.Printers {
		printer.Print(fmt.Sprintf("[Info] %s", message))
	}
}

// LogError ..
func (logger *Logger) LogError(message string) {
	for _, printer := range logger.Printers {
		printer.Print(fmt.Sprintf("[Error] %s", message))
	}
}
