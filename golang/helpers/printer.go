package helpers

import "fmt"

type NoOpPrinter struct{}

func (n NoOpPrinter) Printf(format string, a ...interface{}) {}

// StdOutPrinter is a Printer that writes to standard output.
type StdOutPrinter struct{}

func (p StdOutPrinter) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

type Printer interface {
	Printf(format string, a ...interface{})
}
