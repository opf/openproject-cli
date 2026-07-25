package printer

import "fmt"

// TestingPrinter captures output for assertions. Result collects what would
// go to standard out, ErrResult what would go to standard error.
type TestingPrinter struct {
	Result    string
	ErrResult string
}

func (printer *TestingPrinter) Printf(format string, a ...any) (n int, err error) {
	printer.Result += fmt.Sprintf(format, a...)

	return len(printer.Result), nil
}

func (printer *TestingPrinter) Println(a ...any) (n int, err error) {
	printer.Result += fmt.Sprintln(a...)

	return len(printer.Result), nil
}

func (printer *TestingPrinter) Eprintln(a ...any) (n int, err error) {
	printer.ErrResult += fmt.Sprintln(a...)

	return len(printer.ErrResult), nil
}

func (printer *TestingPrinter) Eprintf(format string, a ...any) (n int, err error) {
	printer.ErrResult += fmt.Sprintf(format, a...)

	return len(printer.ErrResult), nil
}

func (printer *TestingPrinter) Reset() {
	printer.Result = ""
	printer.ErrResult = ""
}
