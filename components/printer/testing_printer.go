package printer

import "fmt"

type TestingPrinter struct {
	Result string
}

func (printer *TestingPrinter) Printf(format string, a ...any) (n int, err error) {
	printer.Result = fmt.Sprintf(format, a...)

	return len(printer.Result), nil
}

func (printer *TestingPrinter) Println(a ...any) (n int, err error) {
	printer.Result = fmt.Sprintln(a...)

	return len(printer.Result), nil
}
