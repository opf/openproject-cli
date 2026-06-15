package printer

import (
	"fmt"
	"os"
)

type ConsolePrinter struct{}

func (printer *ConsolePrinter) Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (printer *ConsolePrinter) Println(a ...any) (n int, err error) {
	return fmt.Println(a...)
}

func (printer *ConsolePrinter) Eprintln(a ...any) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
