package printer

type Printer interface {
	Printf(format string, a ...any) (n int, err error)
	Println(a ...any) (n int, err error)
	// Eprintln writes a line to standard error, keeping diagnostics out of the
	// machine-readable output written to standard out.
	Eprintln(a ...any) (n int, err error)
}

var activePrinter Printer

func Init(printer Printer) {
	activePrinter = printer
}
