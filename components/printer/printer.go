package printer

type Printer interface {
	Printf(format string, a ...any) (n int, err error)
	Println(a ...any) (n int, err error)
}

var activePrinter Printer

func Init(printer Printer) {
	activePrinter = printer
}
