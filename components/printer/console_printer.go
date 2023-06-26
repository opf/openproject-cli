package printer

import "fmt"

type ConsolePrinter struct{}

func (printer *ConsolePrinter) Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (printer *ConsolePrinter) Println(a ...any) (n int, err error) {
	return fmt.Println(a...)
}
