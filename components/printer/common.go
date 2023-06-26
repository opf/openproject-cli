package printer

import "os"

func Info(msg string) {
	activePrinter.Println(msg)
}

func Done() {
	activePrinter.Println(green("DONE"))
}

func Error(err error) {
	activePrinter.Printf("%s Program exited with error: %+v\n", red("[ERROR]"), err)
	os.Exit(-1)
}

func ErrorText(msg string) {
	activePrinter.Printf("%s %s\n", red("[ERROR]"), msg)
	os.Exit(-1)
}

func ResponseError(status int, body []byte) {
	var bodyRepresentation string
	if len(body) >= 256 {
		bodyRepresentation = string(body[:256]) + "\n..."
	} else {
		bodyRepresentation = string(body)
	}

	activePrinter.Printf(
		"%s Bad response from server: (%d)\n\n%s\n",
		red("[ERROR]"),
		status,
		bodyRepresentation,
	)
	os.Exit(-1)
}

func indent(spaces int) (res string) {
	for len(res) < spaces {
		res += " "
	}
	return res
}
